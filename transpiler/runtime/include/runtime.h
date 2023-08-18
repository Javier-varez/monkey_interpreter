#pragma once

#include <builtins.h>
#include <object.h>

namespace runtime {

inline Object rangeExprToArray(const Object start, const Object end) noexcept {
  using std::literals::operator""sv;
  check(start.type == ObjectType::INTEGER && end.type == ObjectType::INTEGER,
        "Cannot construct range expression from arguments of type "sv,
        objectTypeToString(start.type), " and "sv,
        objectTypeToString(end.type));

  return Object::makeArray(
      Array::makeFromRange(start.getInteger(), end.getInteger()));
}

}
