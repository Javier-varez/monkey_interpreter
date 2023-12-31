#pragma once

#include <array.h>
#include <fatal.h>
#include <function.h>

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

class VarArgs;
class HashMap;

struct Object final {
  // Marker type just to make sure nil is represented by the variant
  struct Nil {};

  constexpr static std::array OBJECT_TYPE_NAMES{[]() {
    using std::literals::operator""sv;
    return std::array{
        "NIL"sv,      "INTEGER"sv, "BOOLEAN"sv, "STRING"sv,
        "FUNCTION"sv, "ARRAY"sv,   "VARARGS"sv, "MAP"sv,
    };
  }()};

  enum class Index : size_t {
    NIL,
    INTEGER,
    BOOLEAN,
    STRING,
    FUNCTION,
    ARRAY,
    VARARGS,
    HASH_MAP,
  };

  using Inner = std::variant<Nil, int64_t, bool, std::string, Function, Array,
                             Rc<VarArgs>, Rc<HashMap>>;
  Inner val{Nil{}};

  static inline Object makeInt(const int64_t val) noexcept {
    return Object{
        .val{val},
    };
  }

  static Object makeBool(const bool val) noexcept {
    return Object{
        .val{val},
    };
  }

  static Object makeString(const std::string_view sv) noexcept;
  static Object makeFunction(const Function f) noexcept;
  static Object makeArray(const Array a) noexcept;
  static Object makeVarargs(const VarArgs &v) noexcept;
  static Object makeHashMap(const HashMap &h) noexcept;

  constexpr inline bool is(const Index idx) const noexcept {
    return val.index() == static_cast<size_t>(idx);
  }

  constexpr inline std::string_view type() const noexcept {
    return OBJECT_TYPE_NAMES[val.index()];
  }

  constexpr inline int64_t getInteger() const noexcept {
    return std::get<int64_t>(val);
  }
  constexpr inline bool getBool() const noexcept { return std::get<bool>(val); }
  std::string getString() const noexcept;
  Array getArray() const noexcept;
  VarArgs getVarArgs() const noexcept;
  HashMap getHashMap() const noexcept;

  [[nodiscard]] std::string inspect() const noexcept;

  template <typename... Args>
  Object operator()(const Args &...args) const noexcept;
  Object operator-() const noexcept;
  Object operator!() const noexcept;
  Object operator[](Object index) const noexcept;

  [[nodiscard]] bool equals(const Object &other) const noexcept;
  [[nodiscard]] std::int64_t hash() const noexcept;

  [[nodiscard]] static const Object &nil() noexcept;
};

Object operator+(const Object &lhs, const Object &rhs) noexcept;
Object operator-(const Object &lhs, const Object &rhs) noexcept;
Object operator*(const Object &lhs, const Object &rhs) noexcept;
Object operator/(const Object &lhs, const Object &rhs) noexcept;
Object operator==(const Object &lhs, const Object &rhs) noexcept;
Object operator!=(const Object &lhs, const Object &rhs) noexcept;
Object operator<(const Object &lhs, const Object &rhs) noexcept;
Object operator>(const Object &lhs, const Object &rhs) noexcept;

}  // namespace runtime
