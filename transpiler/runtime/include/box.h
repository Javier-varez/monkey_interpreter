#ifndef _RUNTIME_BOX_H_
#define _RUNTIME_BOX_H_

#include <utility>

namespace runtime {

template <typename T> class Box final {
public:
  template <typename... Args>
  Box(Args &&...args) : mInner{new T{std::forward<Args>(args)...}} {}

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
