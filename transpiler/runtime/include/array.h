#pragma once

#include <vec.h>

namespace runtime {

struct Object;

class Array {
public:
  using Iter = Iterator<const Object>;

  Array() noexcept = default;

  template <typename Arg, typename... Args>
    requires(!Callable<Arg, void, LargeVec<Object>::Pusher>)
  explicit Array(Arg &&arg, Args &&...args) noexcept;

  template <Callable<void, typename LargeVec<Object>::Pusher> C>
  Array(C callable, const size_t sizeHint = 0) noexcept;

  static Array makeFromRange(int64_t start, int64_t end) noexcept;
  static Array makeFromIters(Iter begin, Iter end) noexcept;

  Object operator[](size_t index) const noexcept;

  size_t len() const noexcept;

  Iter begin() const noexcept;
  Iter end() const noexcept;

  Array push(const Object &obj) const noexcept;

private:
  LargeVec<Object> data;
};

template <typename Arg, typename... Args>
  requires(!Callable<Arg, void, LargeVec<Object>::Pusher>)
Array::Array(Arg &&arg, Args &&...args) noexcept
    : data{std::forward<Arg>(arg), std::forward<Args>(args)...} {}

template <Callable<void, typename LargeVec<Object>::Pusher> C>
Array::Array(C callable, const size_t sizeHint) noexcept
    : data{callable, sizeHint} {}

}
