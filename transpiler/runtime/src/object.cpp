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
}  // namespace

[[nodiscard]] std::string Object::inspect() const noexcept {
  return std::visit(Printer{}, val);
}


Object Object::makeInt(const int64_t val) noexcept {
  return Object{
      .type = ObjectType::INTEGER,
      .val{val},
  };
}

Object Object::makeBool(const bool val) noexcept {
  return Object{
      .type = ObjectType::BOOLEAN,
      .val{val},
  };
}

Object Object::makeString(const std::string_view sv) noexcept {
  return Object{
      .type = ObjectType::STRING,
      .val{std::string{sv}},
  };
}

Object Object::makeFunction(const Function f) noexcept {
  return Object{
      .type = ObjectType::FUNCTION,
      .val{f},
  };
}

Object Object::makeArray(const Array a) noexcept {
  return Object{
      .type = ObjectType::ARRAY,
      .val{a},
  };
}

Object Object::makeVarargs(const VarArgs& v) noexcept {
  return Object{
      .type = ObjectType::VARARGS,
      .val{Rc<VarArgs>{Marker<VarArgs>{}, v}},
  };
}

int64_t Object::getInteger() const noexcept {
  using std::literals::operator""sv;
  check(type == ObjectType::INTEGER,
        "Attempted to unwrap integer but object type was `"sv,
        objectTypeToString(type), '`');
  return std::get<int64_t>(val);
}

bool Object::getBool() const noexcept {
  using std::literals::operator""sv;
  check(type == ObjectType::BOOLEAN,
        "Attempted to unwrap bool but object type was `"sv,
        objectTypeToString(type), '`');
  return std::get<bool>(val);
}

std::string Object::getString() const noexcept {
  using std::literals::operator""sv;
  check(type == ObjectType::STRING,
        "Attempted to unwrap string but object type was `"sv,
        objectTypeToString(type), '`');
  return std::get<std::string>(val);
}

Array Object::getArray() const noexcept {
  using std::literals::operator""sv;
  check(type == ObjectType::ARRAY,
        "Attempted to unwrap array but object type was `"sv,
        objectTypeToString(type), '`');
  return std::get<Array>(val);
}

VarArgs Object::getVarArgs() const noexcept {
  using std::literals::operator""sv;
  check(type == ObjectType::VARARGS,
        "Attempted to unwrap varargs but object type was `"sv,
        objectTypeToString(type), '`');
  return *std::get<Rc<VarArgs>>(val);
}

Object Object::operator-() const noexcept {
  using std::literals::operator""sv;
  check(type == ObjectType::INTEGER,
        "Attempted to execute prefix operator '-' on a "sv,
        objectTypeToString(type));
  return Object::makeInt(-getInteger());
}

Object Object::operator!() const noexcept {
  using std::literals::operator""sv;
  check(type == ObjectType::BOOLEAN,
        "Attempted to execute prefix operator '!' on a "sv,
        objectTypeToString(type));
  return Object::makeBool(!getBool());
}

Object Object::operator[](Object index) const noexcept {
  using std::literals::operator""sv;
  if (type == ObjectType::ARRAY) {
    check(index.type == ObjectType::INTEGER,
          "Index to array is not an integer"sv);
    return getArray()[index.getInteger()];
  }

  fatal("Attempted to use index operator on an unsupported object: "sv,
        objectTypeToString(type));
}

Object operator+(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeInt(int64_t{lhs.getInteger() + rhs.getInteger()});
  } else if (lhs.type == ObjectType::STRING && rhs.type == ObjectType::STRING) {
    return Object::makeString(lhs.getString() + rhs.getString());
  }

  fatal("Operator `+` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

Object operator-(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeInt(int64_t{lhs.getInteger() - rhs.getInteger()});
  }

  fatal("Operator `-` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

Object operator*(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeInt(int64_t{lhs.getInteger() * rhs.getInteger()});
  }

  fatal("Operator `*` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

Object operator/(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeInt(int64_t{lhs.getInteger() / rhs.getInteger()});
  }

  fatal("Operator `/` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

Object operator==(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeBool(lhs.getInteger() == rhs.getInteger());
  } else if (lhs.type == ObjectType::BOOLEAN &&
             rhs.type == ObjectType::BOOLEAN) {
    return Object::makeBool(lhs.getBool() == rhs.getBool());
  }

  fatal("Operator `==` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

Object operator!=(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeBool(lhs.getInteger() != rhs.getInteger());
  } else if (lhs.type == ObjectType::BOOLEAN &&
             rhs.type == ObjectType::BOOLEAN) {
    return Object::makeBool(lhs.getBool() != rhs.getBool());
  }

  fatal("Operator `!=` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

Object operator<(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeBool(lhs.getInteger() < rhs.getInteger());
  }

  fatal("Operator `<` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

Object operator>(const Object &lhs, const Object &rhs) noexcept {
  using std::literals::operator""sv;
  if (lhs.type == ObjectType::INTEGER && rhs.type == ObjectType::INTEGER) {
    return Object::makeBool(lhs.getInteger() > rhs.getInteger());
  }

  fatal("Operator `>` is undefined for operands `"sv,
        objectTypeToString(lhs.type), "` and `"sv, objectTypeToString(rhs.type),
        '\n');
}

}
