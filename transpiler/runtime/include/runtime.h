#include <cstdarg>
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
  MAP,
};

struct Object;

template <typename T, T> struct ConstexprLit {};

class Function final {
private:
  struct Callable {
    virtual Object vcall(size_t count, va_list arg) const noexcept = 0;
  };

  template <typename T, size_t NumArgs, bool HasVarArgs>
  struct CallableImpl final : public Callable {
    CallableImpl(T callable) : callable{callable} {}

    Object vcall(size_t count, va_list arg) const noexcept final;

    T callable;
  };

public:
  template <typename T, size_t NumArgs, bool HasVarArgs>
  Function(ConstexprLit<size_t, NumArgs>, ConstexprLit<bool, HasVarArgs>,
           T &&callable)
      : callable{
            std::make_shared<CallableImpl<T, NumArgs, HasVarArgs>>(callable)} {}

  inline Object operator()(size_t count, ...) const noexcept;

private:
  std::shared_ptr<Callable> callable;
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
  case MAP:
    return "map"sv;
  }

  fatal("Invalid object type: "sv, static_cast<int>(type));
}

struct Object final {
  // Marker type just to make sure nil is represented by the variant
  struct Nil {};

  ObjectType type{ObjectType::NIL};
  // TODO: Array object
  // TODO: Map object
  std::variant<Nil, int64_t, bool, std::string, Function> val{Nil{}};

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

  inline static Object makeString(std::string_view sv) noexcept {
    return Object{
        .type = ObjectType::STRING,
        .val{std::string{sv}},
    };
  }

  inline static Object makeFunction(Function f) noexcept {
    return Object{
        .type = ObjectType::FUNCTION,
        .val{f},
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
  };

  [[nodiscard]] inline std::string inspect() const noexcept {
    return std::visit(Printer{}, val);
  }

  template <typename... Args>
  Object operator()(const Args &...args) const noexcept;

  inline Object operator-() const noexcept;

  inline Object operator!() const noexcept;
};

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
  return f(sizeof...(Args), &args...);
}

inline Object Object::operator-() const noexcept {
  check(type == ObjectType::INTEGER, "Attempted to execute prefix operator '-' on a "sv, objectTypeToString(type));
  return Object::makeInt(-getInteger());
}

inline Object Object::operator!() const noexcept {
  check(type == ObjectType::BOOLEAN, "Attempted to execute prefix operator '!' on a "sv, objectTypeToString(type));
  return Object::makeBool(!getBool());
}

Object Function::operator()(size_t count, ...) const noexcept {
  va_list args;
  va_start(args, count);
  const auto result = callable->vcall(count, args);
  va_end(args);
  return result;
}

template <size_t NumArgs, typename C, typename... Args>
auto expandAndCall(va_list arg, C &&callable, Args... args) {
  if constexpr (NumArgs == 0) {
    return callable(std::forward<Args>(args)...);
  } else {
    const auto next = *va_arg(arg, Object *);
    return expandAndCall<NumArgs - 1>(arg, std::forward<C>(callable), args...,
                                      next);
  }
}

template <typename T, size_t NumArgs, bool HasVarArgs>
Object Function::CallableImpl<T, NumArgs, HasVarArgs>::vcall(
    size_t count, va_list arg) const noexcept {
  if constexpr (!HasVarArgs) {
    check(count == NumArgs, "Callable takes "sv, NumArgs, " arguments, but "sv,
          count, " were given");
    return expandAndCall<NumArgs>(arg, callable);
  } else {
    fatal("VarArgs are not implemented yet"sv);
  }
}

template <typename... Args> Object puts(Args &&...args) noexcept {
  const auto print = []<typename T>(T &&arg) { std::cout << arg.inspect(); };
  (print(std::forward<Args>(args)), ...);
  std::cout << '\n';
  return Object{};
}

} // namespace runtime
