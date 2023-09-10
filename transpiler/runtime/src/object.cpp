#include <hash_map.h>
#include <object.h>
#include <var_args.h>

namespace runtime {

namespace {
struct Printer {
  [[nodiscard]] std::string operator()(const Object::Nil &val) noexcept {
    using std::literals::operator""s;
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
    using std::literals::operator""s;
    if (val) {
      return "true"s;
    }

    return "false"s;
  }

  [[nodiscard]] std::string operator()(const Function &val) noexcept {
    using std::literals::operator""s;
    return "<Function>"s;
  }

  [[nodiscard]] std::string operator()(const Array &val) noexcept {
    using std::literals::operator""sv;
    std::ostringstream stream;
    stream << '[';
    bool firstIter = true;
    for (const Object &obj : val) {
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

  [[nodiscard]] std::string operator()(const Rc<VarArgs> &val) noexcept {
    using std::literals::operator""sv;
    std::ostringstream stream;
    stream << "VarArgs["sv;
    bool firstIter = true;
    for (const Object &obj : *val) {
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

  [[nodiscard]] std::string operator()(const Rc<HashMap> &val) noexcept {
    using std::literals::operator""sv;
    std::ostringstream stream;
    stream << "{"sv;
    bool firstIter = true;
    val->forEach([&firstIter, &stream](const Object &k, const Object &v) {
      if (!firstIter) {
        stream << ", "sv;
      } else {
        firstIter = false;
      }
      stream << k.inspect();
      stream << ": ";
      stream << v.inspect();
    });
    stream << '}';
    return stream.str();
  }
};
}  // namespace

[[nodiscard]] std::string Object::inspect() const noexcept {
  return std::visit(Printer{}, val);
}

Object Object::makeString(const std::string_view sv) noexcept {
  return Object{
      .val{std::string{sv}},
  };
}

Object Object::makeFunction(const Function f) noexcept {
  return Object{
      .val{f},
  };
}

Object Object::makeArray(const Array a) noexcept {
  return Object{
      .val{a},
  };
}

Object Object::makeVarargs(const VarArgs &v) noexcept {
  return Object{
      .val{Rc<VarArgs>{Marker<VarArgs>{}, v}},
  };
}

Object Object::makeHashMap(const HashMap &v) noexcept {
  return Object{
      .val{Rc<HashMap>{Marker<HashMap>{}, v}},
  };
}

std::string Object::getString() const noexcept {
  using std::literals::operator""sv;
  check(is(Index::STRING), "Attempted to unwrap string but object type was `"sv,
        type(), '`');
  return std::get<std::string>(val);
}

Array Object::getArray() const noexcept {
  using std::literals::operator""sv;
  check(is(Index::ARRAY), "Attempted to unwrap array but object type was `"sv,
        type(), '`');
  return std::get<Array>(val);
}

VarArgs Object::getVarArgs() const noexcept {
  using std::literals::operator""sv;
  check(is(Index::VARARGS),
        "Attempted to unwrap varargs but object type was `"sv, type(), '`');
  return *std::get<Rc<VarArgs>>(val);
}

HashMap Object::getHashMap() const noexcept {
  using std::literals::operator""sv;
  check(is(Index::HASH_MAP),
        "Attempted to unwrap HashMap but object type was `"sv, type(), '`');
  return *std::get<Rc<HashMap>>(val);
}

Object Object::operator-() const noexcept {
  using std::literals::operator""sv;
  check(is(Index::INTEGER), "Attempted to execute prefix operator '-' on a "sv,
        type());
  return Object::makeInt(-getInteger());
}

Object Object::operator!() const noexcept {
  using std::literals::operator""sv;
  check(is(Index::BOOLEAN), "Attempted to execute prefix operator '!' on a "sv,
        type());
  return Object::makeBool(!getBool());
}

Object Object::operator[](Object index) const noexcept {
  using std::literals::operator""sv;
  if (is(Index::ARRAY)) {
    check(index.is(Index::INTEGER), "Index to array is not an integer: "sv,
          type());
    return getArray()[index.getInteger()];
  } else if (is(Index::HASH_MAP)) {
    return getHashMap()[index];
  }
  fatal("Attempted to use index operator on an unsupported object: "sv, type());
}

Object operator+(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.is(Object::Index::INTEGER) && rhs.is(Object::Index::INTEGER)) {
    return Object::makeInt(int64_t{lhs.getInteger() + rhs.getInteger()});
  } else if (lhs.is(Object::Index::STRING) && rhs.is(Object::Index::STRING)) {
    return Object::makeString(lhs.getString() + rhs.getString());
  }

  fatal("Operator `+` is undefined for operands `"sv, lhs.type(), "` and `"sv,
        rhs.type(), '\n');
}

Object operator-(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.is(Object::Index::INTEGER) && rhs.is(Object::Index::INTEGER)) {
    return Object::makeInt(int64_t{lhs.getInteger() - rhs.getInteger()});
  }

  fatal("Operator `-` is undefined for operands `"sv, lhs.type(), "` and `"sv,
        rhs.type(), '\n');
}

Object operator*(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.is(Object::Index::INTEGER) && rhs.is(Object::Index::INTEGER)) {
    return Object::makeInt(int64_t{lhs.getInteger() * rhs.getInteger()});
  }

  fatal("Operator `*` is undefined for operands `"sv, lhs.type(), "` and `"sv,
        rhs.type(), '\n');
}

Object operator/(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.is(Object::Index::INTEGER) && rhs.is(Object::Index::INTEGER)) {
    return Object::makeInt(int64_t{lhs.getInteger() / rhs.getInteger()});
  }

  fatal("Operator `/` is undefined for operands `"sv, lhs.type(), "` and `"sv,
        rhs.type(), '\n');
}

Object operator==(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.is(Object::Index::INTEGER) && rhs.is(Object::Index::INTEGER)) {
    return Object::makeBool(lhs.getInteger() == rhs.getInteger());
  } else if (lhs.is(Object::Index::BOOLEAN) && rhs.is(Object::Index::BOOLEAN)) {
    return Object::makeBool(lhs.getBool() == rhs.getBool());
  }

  fatal("Operator `==` is undefined for operands `"sv, lhs.type(), "` and `"sv,
        rhs.type(), '\n');
}

Object operator!=(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.is(Object::Index::INTEGER) && rhs.is(Object::Index::INTEGER)) {
    return Object::makeBool(lhs.getInteger() != rhs.getInteger());
  } else if (lhs.is(Object::Index::BOOLEAN) && rhs.is(Object::Index::BOOLEAN)) {
    return Object::makeBool(lhs.getBool() != rhs.getBool());
  }

  fatal("Operator `!=` is undefined for operands `"sv, lhs.type(), "` and `"sv,
        rhs.type(), '\n');
}

Object operator<(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.is(Object::Index::INTEGER) && rhs.is(Object::Index::INTEGER)) {
    return Object::makeBool(lhs.getInteger() < rhs.getInteger());
  }

  fatal("Operator `<` is undefined for operands `"sv, lhs.type(), "` and `"sv,
        rhs.type(), '\n');
}

Object operator>(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.is(Object::Index::INTEGER) && rhs.is(Object::Index::INTEGER)) {
    return Object::makeBool(lhs.getInteger() > rhs.getInteger());
  }

  fatal("Operator `>` is undefined for operands `"sv, lhs.type(), "` and `"sv,
        rhs.type(), '\n');
}

const Object &Object::nil() noexcept {
  const static Object obj{};
  return obj;
}

bool Object::equals(const Object &other) const noexcept {
  return std::visit(
      [&other]<typename T>(const T &lhs) noexcept -> bool {
        if (!std::holds_alternative<T>(other.val)) {
          return false;
        }

        const T &rhs = std::get<T>(other.val);
        if constexpr (std::same_as<T, Nil>) {
          return true;
        } else if constexpr (std::same_as<T, Function>) {
          return false;
        } else if constexpr (std::same_as<T, Array>) {
          return false;
        } else if constexpr (std::same_as<T, Rc<VarArgs>>) {
          return false;
        } else if constexpr (std::same_as<T, Rc<HashMap>>) {
          return false;
        } else {
          return rhs == lhs;
        }
      },
      val);
}

int64_t Object::hash() const noexcept {
  return std::hash<size_t>()(val.index()) ^
         std::visit(
             [this]<typename T>(const T &val) noexcept -> int64_t {
               if constexpr (std::same_as<T, Nil>) {
                 return 0;
               } else if constexpr (std::same_as<T, Function> ||
                                    std::same_as<T, Array> ||
                                    std::same_as<T, Rc<VarArgs>> ||
                                    std::same_as<T, Rc<HashMap>>) {
                 using std::literals::operator""sv;
                 fatal("Cannot hash type: "sv, type());
               } else {
                 return std::hash<T>{}(val);
               }
             },
             val);
}

}  // namespace runtime
