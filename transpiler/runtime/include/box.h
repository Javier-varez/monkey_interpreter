#ifndef _RUNTIME_BOX_H_
#define _RUNTIME_BOX_H_

#include <concepts>
#include <utility>

namespace runtime {

namespace detail {
template <typename U, typename T>
concept SameOrDerived = std::derived_from<U, T> || std::same_as<U, T>;
}

template <typename U> struct Marker final {};

template <typename T> class Box final {
public:
  Box() noexcept : mInner{new T{}} {}

  template <detail::SameOrDerived<T> U, typename... Args>
  Box(Marker<U>, Args &&...args) : mInner{new U{std::forward<Args>(args)...}} {}

  Box(Box &) = delete;
  Box(const Box &) = delete;
  Box(Box &&) = delete;

  Box &operator=(const Box &) = delete;
  Box &operator=(Box &&) = delete;

  constexpr ~Box() noexcept { delete mInner; }

  constexpr T *operator->() noexcept { return mInner; }

  constexpr const T *operator->() const noexcept { return mInner; }

  constexpr T &operator*() noexcept { return *mInner; }

  constexpr const T &operator*() const noexcept { return *mInner; }

private:
  T *mInner;
};

} // namespace runtime

#endif // _RUNTIME_BOX_H_
