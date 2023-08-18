#pragma once

#include <string_view>

#include <fatal.h>

namespace runtime {

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

[[nodiscard]] inline std::string_view
objectTypeToString(const ObjectType type) {
  using std::literals::operator""sv;
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

} // namespace runtime
