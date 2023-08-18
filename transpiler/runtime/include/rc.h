#ifndef _RUNTIME_RC_H_
#define _RUNTIME_RC_H_

#include <cstdint>
#include <cstdlib>
#include <utility>

#include <box.h>

namespace runtime {

/**
 * \brief non-atomic reference count
 */
template <typename T> class Rc {
public:
  Rc() noexcept : mBlock{new Block{}} {}

  template <detail::SameOrDerived<T> U, typename... Args>
  constexpr Rc(Marker<U>, Args &&...args) noexcept
      : mBlock{
            new Block{.elem{Box<T>{Marker<U>{}, std::forward<Args>(args)...}}}} {}

  constexpr Rc(const Rc &other) noexcept : mBlock{other.mBlock} {
    ++mBlock->refCount;
  }

  constexpr Rc(Rc &other) noexcept : mBlock{other.mBlock} {
    ++mBlock->refCount;
  }

  constexpr Rc(Rc &&other) noexcept : mBlock{other.mBlock} {
    ++mBlock->refCount;
  }

  constexpr Rc &operator=(const Rc &other) noexcept {
    if (this != &other) {
      destroy();

      mBlock = other.mBlock;
      ++mBlock->refCount;
    }
    return *this;
  }

  constexpr Rc &operator=(Rc &&other) noexcept {
    if (this != &other) {
      destroy();

      mBlock = other.mBlock;
      ++mBlock->refCount;
    }
    return *this;
  }

  constexpr ~Rc() noexcept { destroy(); }

  constexpr T *operator->() noexcept { return &*mBlock->elem; }

  constexpr const T *operator->() const noexcept { return &*mBlock->elem; }

  constexpr T &operator*() noexcept { return *mBlock->elem; }

  constexpr const T &operator*() const noexcept { return *mBlock->elem; }

private:
  struct Block {
    Box<T> elem;
    size_t refCount{1};
  };

  Block *mBlock;

  constexpr void destroy() {
    --mBlock->refCount;

    if (mBlock->refCount == 0) {
      delete mBlock;
    }
  }
};

} // namespace runtime

#endif // _RUNTIME_RC_H_
