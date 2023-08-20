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
};
} // namespace

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

} // namespace runtime
