#pragma once

#include <object.h>
#include <vec.h>

namespace runtime {

struct Object;

class VarArgs {
public:
  using Iter = Iterator<const Object>;

  VarArgs(Iter begin, Iter end) noexcept;

  size_t len() const noexcept;

  Object operator[](size_t idx) const noexcept;

  Iter begin() const noexcept;
  Iter end() const noexcept;

private:
  Vec<Object> args;
};

} // namespace runtime
