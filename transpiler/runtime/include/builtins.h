#pragma once

#include <object.h>
#include <var_args.h>

#include <function_impl.h>
#include <object_impl.h>

namespace runtime {

template <typename... Args> Object puts(Args &&...args) noexcept {
  const auto print = []<typename T>(T &&arg) { std::cout << arg.inspect(); };

  const auto expandVarArgs = []<typename C, typename T>(C callable, T &&arg) {
    if (arg.is(Object::Index::VARARGS)) {
      for (const Object &inner : arg.getVarArgs()) {
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
  using std::literals::operator""sv;
  check(object.is(Object::Index::VARARGS),
        "Unsupported object passed to toArray: "sv, object.type());

  const VarArgs varargs = object.getVarArgs();
  return object.makeArray(Array::makeFromIters(varargs.begin(), varargs.end()));
}

inline Object len(Object object) noexcept {
  using std::literals::operator""sv;
  check(object.is(Object::Index::ARRAY), "Unsupported object passed to len: "sv,
        object.type());

  const Array arr = object.getArray();
  return Object::makeInt(arr.len());
}

inline Object first(Object object) noexcept {
  using std::literals::operator""sv;
  check(object.is(Object::Index::ARRAY),
        "Unsupported object passed to first: "sv, object.type());

  const Array arr = object.getArray();
  const size_t length = arr.len();
  check(length >= 1,
        "Array does not have any items. Unable to get first item"sv,
        object.type());
  return arr[0];
}

inline Object last(Object object) noexcept {
  using std::literals::operator""sv;
  check(object.is(Object::Index::ARRAY),
        "Unsupported object passed to first: "sv, object.type());

  const Array arr = object.getArray();
  const size_t length = arr.len();
  check(length >= 1, "Array does not have any items. Unable to get last item"sv,
        object.type());
  return arr[length - 1];
}

inline Object rest(Object object) noexcept {
  using std::literals::operator""sv;
  check(object.is(Object::Index::ARRAY),
        "Unsupported object passed to first: "sv, object.type());

  const Array arr = object.getArray();
  const size_t length = arr.len();
  check(length >= 1, "Array does not have any items, rest may not be called"sv,
        object.type());
  return Object::makeArray(Array::makeFromIters(arr.begin() + 1, arr.end()));
}

inline Object push(Object object, Object newObj) noexcept {
  using std::literals::operator""sv;
  check(object.is(Object::Index::ARRAY),
        "Unsupported object passed to first: "sv, object.type());

  const auto &arr = object.getArray();
  auto newArray = arr.push(newObj);
  return Object::makeArray(newArray);
}

} // namespace runtime
