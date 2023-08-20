#pragma once

#include <builtins.h>
#include <object.h>

namespace runtime {

inline Object rangeExprToArray(const Object start, const Object end) noexcept {
  using std::literals::operator""sv;
  check(start.is(Object::Index::INTEGER) && end.is(Object::Index::INTEGER),
        "Cannot construct range expression from arguments of type "sv,
        start.type(), " and "sv, end.type());

  return Object::makeArray(
      Array::makeFromRange(start.getInteger(), end.getInteger()));
}

} // namespace runtime
