#pragma once

#include <cstdint>
#include <cstdlib>
#include <functional>
#include <iostream>
#include <memory>
#include <sstream>
#include <string>
#include <string_view>
#include <variant>

#include <object_type.h>
#include <function.h>
#include <array.h>
#include <fatal.h>

namespace runtime {

class VarArgs;

struct Object final {
  // Marker type just to make sure nil is represented by the variant
  struct Nil {};

  ObjectType type{ObjectType::NIL};

  // TODO: Map object
  std::variant<Nil, int64_t, bool, std::string, Function, Array, Rc<VarArgs>> val{
      Nil{}};

  static Object makeInt(const int64_t val) noexcept;
  static Object makeBool(const bool val) noexcept;
  static Object makeString(const std::string_view sv) noexcept ;
  static Object makeFunction(const Function f) noexcept ;
  static Object makeArray(const Array a) noexcept;
  static Object makeVarargs(const VarArgs& v) noexcept;

  int64_t getInteger() const noexcept;
  bool getBool() const noexcept;
  std::string getString() const noexcept;
  Array getArray() const noexcept;
  VarArgs getVarArgs() const noexcept;

  [[nodiscard]] std::string inspect() const noexcept;

  template <typename... Args>
  Object operator()(const Args &...args) const noexcept;
  Object operator-() const noexcept;
  Object operator!() const noexcept;
  Object operator[](Object index) const noexcept;
};

Object operator+(const Object &lhs, const Object &rhs) noexcept;
Object operator-(const Object &lhs, const Object &rhs) noexcept;
Object operator*(const Object &lhs, const Object &rhs) noexcept;
Object operator/(const Object &lhs, const Object &rhs) noexcept;
Object operator==(const Object &lhs, const Object &rhs) noexcept;
Object operator!=(const Object &lhs, const Object &rhs) noexcept;
Object operator<(const Object &lhs, const Object &rhs) noexcept;
Object operator>(const Object &lhs, const Object &rhs) noexcept;

}
