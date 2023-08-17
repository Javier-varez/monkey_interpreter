#ifndef _RUNTIME_RC_H_
#define _RUNTIME_RC_H_

#include <cstdint>
#include <cstdlib>
#include <utility>

namespace runtime {

/**
 * \brief non-atomic reference count
 */
template <typename T> class Rc {
public:
  template <typename... Args>
  constexpr Rc(Args &&...args) noexcept
      : mBlock{
            new Block{.elem{T{std::forward<Args>(args)...}}, .refCount = 1}} {}

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

  constexpr T *operator->() noexcept { return &mBlock->elem; }

  constexpr const T *operator->() const noexcept { return &mBlock->elem; }

  constexpr T &operator*() noexcept { return mBlock->elem; }

  constexpr const T &operator*() const noexcept { return mBlock->elem; }

private:
  struct Block {
    T elem;
    size_t refCount;
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
