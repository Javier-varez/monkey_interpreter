#ifndef _MONKEY_RUNTIME_H_
#define _MONKEY_RUNTIME_H_

#include <cstdint>
#include <cstdlib>
#include <functional>
#include <iostream>
#include <memory>
#include <sstream>
#include <string>
#include <string_view>
#include <variant>

namespace runtime {

using std::literals::operator""s;
using std::literals::operator""sv;

template <typename... Args>
[[noreturn]] inline void fatal(Args &&...args) noexcept {
  const auto print = []<typename T>(T &&arg) {
    std::cout << std::forward<T>(arg);
    return true;
  };
  std::cout << "Failed assertion: "sv;

  (print(std::forward<Args>(args)) && ...);
  std::exit(-1);
}

template <typename... Args>
constexpr inline void check(const bool condition, Args &&...args) noexcept {
  if (!condition) {
    fatal(std::forward<Args>(args)...);
  }
}

enum class ObjectType {
  NIL,
  INTEGER,
  BOOLEAN,
  STRING,
  FUNCTION,
  ARRAY,
  VARARGS,
  MAP,
};

struct Object;

template <typename T, T> struct ConstexprLit {};

class FnArgs {
public:
  template<std::same_as<Object>... Args>
  explicit FnArgs(const Args... args) noexcept;

  inline size_t len() const noexcept;

  inline Object operator[](size_t idx) const noexcept;

  using Iterator = std::vector<Object>::const_iterator;

  inline Iterator begin() const noexcept;
  inline Iterator end() const noexcept;

private:
  std::shared_ptr<std::vector<Object>> args;
};

class VarArgs {
public:
  using Iterator = std::vector<Object>::const_iterator;

  inline VarArgs(Iterator begin, Iterator end) noexcept;

  inline size_t len() const noexcept;

  inline Object operator[](size_t idx) const noexcept;

  inline Iterator begin() const noexcept;
  inline Iterator end() const noexcept;

private:
  std::shared_ptr<std::vector<Object>> args;
};

class Function final {
private:
  struct Callable {
    virtual Object call(const FnArgs args) const noexcept = 0;
  };

  template <typename T, size_t NumArgs, bool HasVarArgs>
  struct CallableImpl final : public Callable {
    CallableImpl(T callable) : callable{callable} {}

    Object call(const FnArgs args) const noexcept final;

    T callable;
  };

public:
  template <typename T, size_t NumArgs, bool HasVarArgs>
  Function(ConstexprLit<size_t, NumArgs>, ConstexprLit<bool, HasVarArgs>,
           T &&callable)
      : callable{
            std::make_shared<CallableImpl<T, NumArgs, HasVarArgs>>(callable)} {}

  inline Object operator()(const FnArgs args) const noexcept;

private:
  std::shared_ptr<Callable> callable;
};

class Array {
public:
  using Iterator = std::vector<Object>::const_iterator;

  inline Array() noexcept;

  template<typename... Args>
  explicit Array(Args&&...args) noexcept;

  inline static Array makeFromRange(int64_t start, int64_t end)noexcept;
  inline static Array makeFromIters(Iterator begin, Iterator end)noexcept;

  inline Object operator[](size_t index) const noexcept;

  inline size_t len() const noexcept;

  inline Iterator begin() const noexcept;
  inline Iterator end() const noexcept;

  inline void push(Object obj) noexcept;

private:
  std::shared_ptr<std::vector<Object>> data;
};

[[nodiscard]] inline std::string_view
objectTypeToString(const ObjectType type) {
  using enum ObjectType;
  switch (type) {
  case NIL:
    return "nil"sv;
  case INTEGER:
    return "integer"sv;
  case BOOLEAN:
    return "boolean"sv;
  case STRING:
    return "string"sv;
  case FUNCTION:
    return "function"sv;
  case ARRAY:
    return "array"sv;
  case VARARGS:
    return "varargs"sv;
  case MAP:
    return "map"sv;
  }

  fatal("Invalid object type: "sv, static_cast<int>(type));
}

struct Object final {
  // Marker type just to make sure nil is represented by the variant
  struct Nil {};

  ObjectType type{ObjectType::NIL};

  // TODO: Map object
  std::variant<Nil, int64_t, bool, std::string, Function, Array, VarArgs> val{Nil{}};

  inline static Object makeInt(const int64_t val) noexcept {
    return Object{
        .type = ObjectType::INTEGER,
        .val{val},
    };
  }

  inline static Object makeBool(const bool val) noexcept {
    return Object{
        .type = ObjectType::BOOLEAN,
        .val{val},
    };
  }

  inline static Object makeString(const std::string_view sv) noexcept {
    return Object{
        .type = ObjectType::STRING,
        .val{std::string{sv}},
    };
  }

  inline static Object makeFunction(const Function f) noexcept {
    return Object{
        .type = ObjectType::FUNCTION,
        .val{f},
    };
  }

  inline static Object makeArray(const Array a) noexcept {
    return Object{
        .type = ObjectType::ARRAY,
        .val{a},
    };
  }

  inline static Object makeVarargs(const VarArgs v) noexcept {
    return Object{
        .type = ObjectType::VARARGS,
        .val{v},
    };
  }

  inline int64_t getInteger() const noexcept {
    check(type == ObjectType::INTEGER,
          "Attempted to unwrap integer but object type was `"sv,
          objectTypeToString(type), '`');
    return std::get<int64_t>(val);
  }

  inline bool getBool() const noexcept {
    check(type == ObjectType::BOOLEAN,
          "Attempted to unwrap bool but object type was `"sv,
          objectTypeToString(type), '`');
    return std::get<bool>(val);
  }

  inline std::string getString() const noexcept {
    check(type == ObjectType::STRING,
          "Attempted to unwrap string but object type was `"sv,
          objectTypeToString(type), '`');
    return std::get<std::string>(val);
  }

  inline Array getArray() const noexcept {
    check(type == ObjectType::ARRAY,
          "Attempted to unwrap array but object type was `"sv,
          objectTypeToString(type), '`');
    return std::get<Array>(val);
  }

  inline VarArgs getVarArgs() const noexcept {
    check(type == ObjectType::VARARGS,
          "Attempted to unwrap varargs but object type was `"sv,
          objectTypeToString(type), '`');
    return std::get<VarArgs>(val);
  }

  [[nodiscard]] inline std::string inspect() const noexcept;

  struct Printer {
    [[nodiscard]] std::string operator()(const Nil &val) noexcept {
      return "nil"s;
    }

    [[nodiscard]] std::string operator()(const std::string &val) noexcept {
      return val;
    }

    [[nodiscard]] std::string operator()(const int64_t val) noexcept {
      std::ostringstream stream;
      stream << val;
      return stream.str();
    }

    [[nodiscard]] std::string operator()(const bool val) noexcept {
      if (val) {
        return "true"s;
      }

      return "false"s;
    }

    [[nodiscard]] std::string operator()(const Function &val) noexcept {
      return "<Function>"s;
    }

    [[nodiscard]] std::string operator()(const Array &val) noexcept {
      std::ostringstream stream;
      stream << '[';
      bool firstIter = true;
      for (const Object& obj : val) {
        if (!firstIter) {
          stream << ", "sv;
        } else {
          firstIter = false;
        }
        stream << obj.inspect();
      }
      stream << ']';
      return stream.str();
    }

    [[nodiscard]] std::string operator()(const VarArgs &val) noexcept {
      std::ostringstream stream;
      stream << "VarArgs["sv;
      bool firstIter = true;
      for (const Object& obj : val) {
        if (!firstIter) {
          stream << ", "sv;
        } else {
          firstIter = false;
        }
        stream << obj.inspect();
      }
      stream << ']';
      return stream.str();
    }
  };

  template <typename... Args>
  Object operator()(const Args &...args) const noexcept;

  inline Object operator-() const noexcept;

  inline Object operator!() const noexcept;

  inline Object operator[](Object index) const noexcept;
};

[[nodiscard]] inline std::string Object::inspect() const noexcept {
  return std::visit(Printer{}, val);
}

inline Object operator+(const Object &lhs, const Object &rhs) noexcept {
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeInt(int64_t{lhs.getInteger() + rhs.getInteger()});
  } else if (lhs.type == ObjectType::STRING && rhs.type == ObjectType::STRING) {
    return Object::makeString(lhs.getString() + rhs.getString());
  }

  fatal("Operator `+` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

inline Object operator-(const Object &lhs, const Object &rhs) noexcept {
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeInt(int64_t{lhs.getInteger() - rhs.getInteger()});
  }

  fatal("Operator `-` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

inline Object operator*(const Object &lhs, const Object &rhs) noexcept {
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeInt(int64_t{lhs.getInteger() * rhs.getInteger()});
  }

  fatal("Operator `*` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

inline Object operator/(const Object &lhs, const Object &rhs) noexcept {
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeInt(int64_t{lhs.getInteger() / rhs.getInteger()});
  }

  fatal("Operator `/` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

inline Object operator==(const Object &lhs, const Object &rhs) noexcept {
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeBool(lhs.getInteger() == rhs.getInteger());
  } else if (lhs.type == ObjectType::BOOLEAN && rhs.type == ObjectType::BOOLEAN) {
    return Object::makeBool(lhs.getBool() == rhs.getBool());
  }

  fatal("Operator `==` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

inline Object operator!=(const Object &lhs, const Object &rhs) noexcept {
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeBool(lhs.getInteger() != rhs.getInteger());
  } else if (lhs.type == ObjectType::BOOLEAN && rhs.type == ObjectType::BOOLEAN) {
    return Object::makeBool(lhs.getBool() != rhs.getBool());
  }

  fatal("Operator `!=` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

inline Object operator<(const Object &lhs, const Object &rhs) noexcept {
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeBool(lhs.getInteger() < rhs.getInteger());
  }

  fatal("Operator `<` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

inline Object operator>(const Object &lhs, const Object &rhs) noexcept {
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeBool(lhs.getInteger() > rhs.getInteger());
  }

  fatal("Operator `>` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

template <typename... Args>
inline Object Object::operator()(const Args &...args) const noexcept {
  const Function f = std::get<Function>(val);
  return f(FnArgs{args...});
}

inline Object Object::operator-() const noexcept {
  check(type == ObjectType::INTEGER, "Attempted to execute prefix operator '-' on a "sv, objectTypeToString(type));
  return Object::makeInt(-getInteger());
}

inline Object Object::operator!() const noexcept {
  check(type == ObjectType::BOOLEAN, "Attempted to execute prefix operator '!' on a "sv, objectTypeToString(type));
  return Object::makeBool(!getBool());
}

inline Object Object::operator[](Object index) const noexcept {
  if (type == ObjectType::ARRAY) {
    check(index.type == ObjectType::INTEGER, "Index to array is not an integer"sv);
    return getArray()[index.getInteger()];
  }

  fatal("Attempted to use index operator on an unsupported object: "sv, objectTypeToString(type));
}

Object Function::operator()(const FnArgs args) const noexcept {
  const auto result = callable->call(args);
  return result;
}

template <size_t NumArgs, typename C, typename... Args>
auto expandAndCall(const FnArgs::Iterator argIter, C &&callable, Args... args) {
  if constexpr (NumArgs == 0) {
    return callable(std::forward<Args>(args)...);
  } else {
    const auto nextIter = argIter+1;
    return expandAndCall<NumArgs - 1>(nextIter, std::forward<C>(callable), args...,
                                      *argIter);
  }
}

template <size_t NumArgs, typename C, typename... Args>
auto expandAndCallWithVarArgs(const FnArgs::Iterator argIter, const FnArgs::Iterator argIterEnd, C &&callable, Args... args) {
  if constexpr (NumArgs == 0) {
    const Object varArgs {Object::makeVarargs(VarArgs{argIter, argIterEnd})};
    return callable(std::forward<Args>(args)..., varArgs);
  } else {
    const auto nextIter = argIter+1;
    return expandAndCallWithVarArgs<NumArgs - 1>(nextIter, argIterEnd, std::forward<C>(callable), args...,
                                                 *argIter);
  }
}

template <typename T, size_t NumArgs, bool HasVarArgs>
Object Function::CallableImpl<T, NumArgs, HasVarArgs>::call(const FnArgs args) const noexcept {
  if constexpr (!HasVarArgs) {
    check(args.len() == NumArgs, "Callable takes "sv, NumArgs, " arguments, but "sv,
          args.len(), " were given");
    return expandAndCall<NumArgs>(args.begin(), callable);
  } else {
    check(args.len() >= NumArgs, "Callable takes at least "sv, NumArgs, " arguments, but only "sv,
          args.len(), " were given");
    return expandAndCallWithVarArgs<NumArgs>(args.begin(), args.end(), callable);
  }
}

Array::Array() noexcept : data{std::make_shared<std::vector<Object>>()} {}

template<typename... Args>
Array::Array(Args&&...args) noexcept : data{std::make_shared<std::vector<Object>>(std::vector<Object>{args...})} {}

Object Array::operator[](size_t index) const noexcept {
  check(index < len(), "Out of bounds access to array.");
  return (*data)[index];
}

size_t Array::len() const noexcept {
  return (*data).size();
}

Array::Iterator Array::begin() const noexcept{
  return (*data).cbegin();
}

Array::Iterator Array::end() const noexcept{
  return (*data).cend();
}

Array Array::makeFromRange(int64_t start, int64_t end) noexcept {
  Array a{};
  const auto abs = [](auto arg) {
    if (arg < 0) return -arg;
    return arg;
  };

  a.data->reserve(abs(end - start));
  if (start > end) {
    for (int64_t i = start; i > end; i--) {
      a.data->push_back(Object::makeInt(i));
    }
  } else {
    for (int64_t i = start; i < end; i++) {
      a.data->push_back(Object::makeInt(i));
    }
  }

  return a;
}

Array Array::makeFromIters(Iterator begin, Iterator end)noexcept {
  Array a{};
  a.data->reserve(end - begin);

  auto next = begin;
  while (next < end) {
    a.data->push_back(*next);
    next++;
  }

  return a;
}

void Array::push(const Object newObj) noexcept {
  data->push_back(newObj);
}

template<std::same_as<Object>... Args>
FnArgs::FnArgs(const Args... args) noexcept : args(std::make_shared<std::vector<Object>>()){
  const auto handleArg = [this]<typename T>(const T arg) {
    static_assert(std::same_as<T, Object>, "Invalid arg type");
    if (arg.type == ObjectType::VARARGS) {
      // Varargs are unwrapped here in the call site
      for (const Object &inner : arg.getVarArgs()) {
        this->args->push_back(inner);
      }
    } else {
      this->args->push_back(arg);
    }
  };

  (handleArg(args), ...);
}

size_t FnArgs::len() const noexcept {
  return args->size();
}

Object FnArgs::operator[](size_t idx) const noexcept {
  check(idx < args->size(), "Out of bounds index to FnArgs object"sv);
  return (*args)[idx];
}

FnArgs::Iterator FnArgs::begin() const noexcept {
  return args->cbegin();
}
FnArgs::Iterator FnArgs::end() const noexcept {
  return args->cend();
}

inline VarArgs::VarArgs(const Iterator begin, const Iterator end) noexcept : args(std::make_shared<std::vector<Object>>()){
  Iterator next= begin;
  while (next < end) {
    args->push_back(*next);
    next++;
  }
}

inline size_t VarArgs::len() const noexcept {
  return args->size();
}

inline Object VarArgs::operator[](size_t idx) const noexcept {
  return (*args)[idx];
}

inline VarArgs::Iterator VarArgs::begin() const noexcept {
  return args->cbegin();
}

inline VarArgs::Iterator VarArgs::end() const noexcept{
  return args->cend();
}

inline Object rangeExprToArray(const Object start, const Object end) noexcept {
  check(start.type == ObjectType::INTEGER &&end.type == ObjectType::INTEGER,
        "Cannot construct range expression from arguments of type "sv, objectTypeToString(start.type),
        " and "sv, objectTypeToString(end.type));

  return Object::makeArray(Array::makeFromRange(start.getInteger(), end.getInteger()));
}

template <typename... Args> Object puts(Args &&...args) noexcept {
  const auto print = []<typename T>(T &&arg) {
    std::cout << arg.inspect();
  };

  const auto expandVarArgs = []<typename C, typename T>(C callable, T&& arg) {
    if (arg.type == ObjectType::VARARGS) {
      for (const Object& inner: arg.getVarArgs()) {
        callable(inner);
      }
    } else {
        callable(arg);
    }
  };

  (expandVarArgs(print, std::forward<Args>(args)), ...);
  std::cout << '\n';
  return Object{};
}

inline Object toArray(Object object) noexcept {
  check(object.type == ObjectType::VARARGS, "Unsupported object passed to toArray: "sv, objectTypeToString(object.type));

  const VarArgs varargs = object.getVarArgs();
  return object.makeArray(Array::makeFromIters(varargs.begin(), varargs.end()));
}

inline Object len(Object object) noexcept {
  check(object.type == ObjectType::ARRAY, "Unsupported object passed to len: "sv, objectTypeToString(object.type));

  const Array arr = object.getArray();
  return Object::makeInt(arr.len());
}

inline Object first(Object object) noexcept {
  check(object.type == ObjectType::ARRAY, "Unsupported object passed to first: "sv, objectTypeToString(object.type));

  const Array arr = object.getArray();
  const size_t length = arr.len();
  check(length >= 1, "Array does not have any items. Unable to get first item"sv, objectTypeToString(object.type));
  return arr[0];
}

inline Object last(Object object) noexcept {
  check(object.type == ObjectType::ARRAY, "Unsupported object passed to first: "sv, objectTypeToString(object.type));

  const Array arr = object.getArray();
  const size_t length = arr.len();
  check(length >= 1, "Array does not have any items. Unable to get last item"sv, objectTypeToString(object.type));
  return arr[length-1];
}

inline Object rest(Object object) noexcept {
  check(object.type == ObjectType::ARRAY, "Unsupported object passed to first: "sv, objectTypeToString(object.type));

  const Array arr = object.getArray();
  const size_t length = arr.len();
  check(length >= 1, "Array does not have any items, rest may not be called"sv, objectTypeToString(object.type));
  return Object::makeArray(Array::makeFromIters(arr.begin()+1, arr.end()));
}

inline Object push(Object object, Object newObj) noexcept {
  check(object.type == ObjectType::ARRAY, "Unsupported object passed to first: "sv, objectTypeToString(object.type));

  const auto arr = object.getArray();
  auto newArray = Array::makeFromIters(arr.begin(), arr.end());
  newArray.push(newObj);
  return Object::makeArray(newArray);
}

} // namespace runtime

#endif  // _MONKEY_RUNTIME_H_
