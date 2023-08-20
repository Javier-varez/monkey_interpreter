#ifndef _RUNTIME_VECTOR_H_
#define _RUNTIME_VECTOR_H_

#include <array>
#include <functional>
#include <type_traits>
#include <variant>

#include <callable.h>
#include <rc.h>

namespace runtime {

/**
 * \brief Memory iterator
 */
template <typename T> class Iterator {
public:
  constexpr Iterator(T *ptr) noexcept : mPtr(ptr) {}

  constexpr Iterator(const Iterator &) noexcept = default;
  constexpr Iterator(Iterator &&) noexcept = default;
  constexpr Iterator &operator=(const Iterator &) noexcept = default;
  constexpr Iterator &operator=(Iterator &&) noexcept = default;

  constexpr Iterator &operator++() noexcept {
    mPtr++;
    return *this;
  }

  constexpr Iterator operator++(int) noexcept {
    const Iterator copy{*this};
    mPtr++;
    return copy;
  }

  constexpr Iterator operator+(size_t val) const noexcept { return mPtr + val; }

  constexpr T *operator->() const noexcept { return mPtr; }

  constexpr T &operator*() const noexcept { return *mPtr; }

  constexpr ssize_t operator-(const Iterator other) const noexcept {
    return mPtr - other.mPtr;
  }

  constexpr auto operator<=>(const Iterator &) const = default;

private:
  T *mPtr;
};

template <typename T> struct Span {
  Iterator<const T> begin;
  Iterator<const T> end;
};

template <typename T, size_t SMALL_VEC_NUM_ELEMS = 10> class Vec;

template <typename T, size_t NUM_ELEMS> class SmallVec final {
public:
  struct Pusher {
    SmallVec &vec;

    constexpr void push(const T &item) const noexcept {
      std::construct_at(reinterpret_cast<T *>(&vec.mStorage[vec.mUsedSize]),
                        item);
      vec.mUsedSize += 1;
    }

    constexpr void push(T &&item) const noexcept {
      std::construct_at(reinterpret_cast<T *>(&vec.mStorage[vec.mUsedSize]),
                        std::move(item));
      vec.mUsedSize += 1;
    }
  };

  template <typename... Args> constexpr SmallVec(Args &&...args) noexcept {
    const Pusher pusher{.vec{*this}};

    const auto handleArg = [pusher]<typename U>(U &&arg) {
      if constexpr (std::same_as<std::remove_cv_t<
                                     std::remove_reference_t<std::decay_t<U>>>,
                                 Span<const T>>) {
        auto current = arg.begin;
        while (current < arg.end) {
          pusher.push(*current);
          ++current;
        }
      } else {
        pusher.push(std::forward<U>(arg));
      }
    };

    (handleArg(std::forward<Args>(args)), ...);
  }

  template <Callable<void, Pusher> C>
  constexpr SmallVec(C callalble, const size_t sizeHint = 0) noexcept {
    callalble(Pusher{.vec{*this}});
  }

  constexpr SmallVec(const SmallVec &other) {
    for (size_t i = 0; i < other.size(); i++) {
      std::construct_at(reinterpret_cast<T *>(&mStorage[mUsedSize]), other[i]);
      mUsedSize += 1;
    }
  }

  constexpr SmallVec(SmallVec &other) {
    for (size_t i = 0; i < other.size(); i++) {
      std::construct_at(reinterpret_cast<T *>(&mStorage[mUsedSize]), other[i]);
      mUsedSize += 1;
    }
  }

  constexpr SmallVec(SmallVec &&other) {
    for (size_t i = 0; i < other.size(); i++) {
      std::construct_at(reinterpret_cast<T *>(&mStorage[mUsedSize]), other[i]);
      mUsedSize += 1;
    }
  }

  constexpr SmallVec &operator=(const SmallVec &other) {
    if (this != &other) {
      for (size_t i = 0; i < other.size(); i++) {
        std::construct_at(reinterpret_cast<T *>(&mStorage[mUsedSize]),
                          other[i]);
        mUsedSize += 1;
      }
    }
    return *this;
  }

  constexpr SmallVec &operator=(SmallVec &&other) {
    if (this != &other) {
      for (size_t i = 0; i < other.size(); i++) {
        std::construct_at(reinterpret_cast<T *>(&mStorage[mUsedSize]),
                          other[i]);
        mUsedSize += 1;
      }
    }
    return *this;
  }

  constexpr ~SmallVec() noexcept {
    for (size_t i = 0; i < mUsedSize; ++i) {
      T &elem = *std::launder(reinterpret_cast<T *>(&mStorage[i]));
      elem.~T();
    }
  }

  constexpr Iterator<const T> begin() const noexcept {
    return Iterator{std::launder(reinterpret_cast<const T *>(&mStorage[0]))};
  }

  constexpr Iterator<const T> end() const noexcept {
    return Iterator{
        std::launder(reinterpret_cast<const T *>(&mStorage[mUsedSize]))};
  }

  constexpr const T &operator[](const size_t index) const noexcept {
    return *std::launder(reinterpret_cast<const T *>(&mStorage[index]));
  }

  constexpr size_t capacity() const noexcept { return NUM_ELEMS; }

  constexpr size_t size() const noexcept { return mUsedSize; }

private:
  using Storage = std::aligned_storage_t<sizeof(T), alignof(T)>;
  std::array<Storage, NUM_ELEMS> mStorage; // Uninitialized for performance
  size_t mUsedSize{};
};

namespace detail {

template <typename T, typename U, typename... Rest>
constexpr static size_t countElems(U &&current, Rest &&...rest) {
  size_t count = 1;
  if constexpr (std::same_as<
                    std::remove_cv_t<std::remove_reference_t<std::decay_t<U>>>,
                    Span<const T>>) {
    count = current.end - current.begin;
  }

  if constexpr (sizeof...(rest) == 0) {
    return count;
  } else {
    return count + countElems<T>(std::forward<Rest>(rest)...);
  }
}

} // namespace detail

template <typename T> class LargeVec final {
public:
  struct Pusher final {
    LargeVec &vec;

    constexpr void push(const T &item) const noexcept {
      vec.mInner->push_back(item);
    }

    void push(T &&item) const noexcept {
      vec.mInner->push_back(std::move(item));
    }
  };

  template <typename... Args> constexpr LargeVec(Args &&...args) noexcept {
    const size_t numElems = detail::countElems<T>(std::forward<Args>(args)...);
    mInner->reserve(numElems);

    const auto handleArg = [this]<typename U>(U &&arg) {
      if constexpr (std::same_as<std::remove_cv_t<
                                     std::remove_reference_t<std::decay_t<U>>>,
                                 Span<const T>>) {
        Iterator<const T> current = arg.begin;
        while (current < arg.end) {
          mInner->push_back(*current);
          ++current;
        }
      } else {
        mInner->push_back(std::forward<U>(arg));
      }
    };

    (handleArg(std::forward<Args>(args)), ...);
  }

  template <Callable<void, Pusher> C>
  constexpr LargeVec(C callable, size_t sizeHint = 0) noexcept {
    if (sizeHint != 0) {
      mInner->reserve(sizeHint);
    }

    // Allow caller to push data during construction
    callable(Pusher{.vec{*this}});
  }

  constexpr Iterator<const T> begin() const noexcept {
    return Iterator{&*mInner->begin()};
  }

  constexpr Iterator<const T> end() const noexcept {
    return Iterator{&*mInner->end()};
  }

  constexpr const T &operator[](const size_t index) const noexcept {
    return (*mInner)[index];
  }

  constexpr size_t capacity() const noexcept { return mInner->capacity(); }

  constexpr size_t size() const noexcept { return mInner->size(); }

  template <typename... Args>
  constexpr LargeVec copyAppend(Args &&...args) const noexcept {
    const Span<const T> currElemsSpan{.begin = begin(), .end = end()};
    return LargeVec{currElemsSpan, std::forward<Args>(args)...};
  }

private:
  Rc<std::vector<T>> mInner;
};

/**
 * \brief Immutable vector implementation
 */
template <typename T, size_t SMALL_VEC_NUM_ELEMS> class Vec {

  using Svec = SmallVec<T, SMALL_VEC_NUM_ELEMS>;
  using Lvec = LargeVec<T>;

public:
  constexpr Vec() noexcept = default;

  template <typename C>
    requires(Callable<C, void, typename Svec::Pusher> &&
             Callable<C, void, typename Lvec::Pusher>)
  constexpr Vec(C callable, const size_t sizeHint = 0) noexcept {
    if (sizeHint == 0 || sizeHint > SMALL_VEC_NUM_ELEMS) {
      mInner.template emplace<Lvec>(callable, sizeHint);
    } else {
      mInner.template emplace<Svec>(callable, sizeHint);
    }
  }

  template <typename... Args>
  constexpr static Vec makeVec(Args &&...args) noexcept {
    const size_t count = detail::countElems<T>(std::forward<Args>(args)...);
    if (count > SMALL_VEC_NUM_ELEMS) {
      // Construct new refcounted
      return Vec{std::in_place_type<Lvec>, std::forward<Args>(args)...};
    }
    return Vec{std::in_place_type<Svec>, std::forward<Args>(args)...};
  }

  constexpr Vec(const Vec &) = default;
  constexpr Vec(Vec &&) = default;
  constexpr Vec &operator=(const Vec &) = default;
  constexpr Vec &operator=(Vec &&) = default;

  constexpr Iterator<const T> begin() const noexcept {
    return std::visit([]<typename I>(const I &elem) { return elem.begin(); },
                      mInner);
  }

  constexpr Iterator<const T> end() const noexcept {
    return std::visit([]<typename I>(const I &elem) { return elem.end(); },
                      mInner);
  }

  constexpr const T &operator[](size_t index) const noexcept {
    return std::visit(
        [index]<typename I>(const I &elem) -> const T & { return elem[index]; },
        mInner);
  }

  constexpr size_t size() const noexcept {
    return std::visit([]<typename I>(const I &elem) { return elem.size(); },
                      mInner);
  }

  constexpr size_t capacity() const noexcept {
    return std::visit([]<typename I>(const I &elem) { return elem.capacity(); },
                      mInner);
  }

  template <typename... Args>
  constexpr Vec copyAppend(Args &&...args) const noexcept {
    // Decide what type of object to create
    const size_t count =
        detail::countElems<T>(std::forward<Args>(args)...) + size();
    const Span<const T> currElemsSpan{.begin = begin(), .end = end()};

    if (count > SMALL_VEC_NUM_ELEMS) {
      // Construct new refcounted
      return Vec{std::in_place_type<Lvec>, currElemsSpan,
                 std::forward<Args>(args)...};
    }

    return Vec{std::in_place_type<Svec>, currElemsSpan,
               std::forward<Args>(args)...};
  }

  constexpr bool isSmallVec() const noexcept { return mInner.index() == 0; }

private:
  std::variant<Svec, Lvec> mInner;

  template <typename... Args>
  constexpr Vec(std::in_place_type_t<Svec>, Args &&...args) noexcept
      : mInner(std::in_place_type<Svec>, std::forward<Args>(args)...) {}

  template <typename... Args>
  constexpr Vec(std::in_place_type_t<Lvec>, Args &&...args) noexcept
      : mInner(std::in_place_type<Lvec>, std::forward<Args>(args)...) {}
};

} // namespace runtime

#endif // _RUNTIME_VECTOR_H_
