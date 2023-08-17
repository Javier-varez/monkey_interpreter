#ifndef _RUNTIME_VECTOR_H_
#define _RUNTIME_VECTOR_H_

#include <array>
#include <type_traits>

#include <rc.h>

namespace runtime {

/**
 * \brief Memory iterator
 */
template <typename T> class Iterator {
public:
  Iterator(T *ptr) noexcept : mPtr(ptr) {}

  Iterator(const Iterator &) noexcept = default;
  Iterator(Iterator &&) noexcept = default;
  Iterator &operator=(const Iterator &) noexcept = default;
  Iterator &operator=(Iterator &&) noexcept = default;

  constexpr Iterator &operator++() noexcept {
    mPtr++;
    return *this;
  }

  constexpr Iterator &operator++(int) noexcept {
    const Iterator copy = *this;
    mPtr++;
    return copy;
  }

  constexpr Iterator operator+(size_t val) const noexcept {
    return mPtr + val;
  }

  constexpr T *operator->() const noexcept { return mPtr; }

  constexpr T &operator*() const noexcept { return *mPtr; }

  constexpr ssize_t operator-(const Iterator other) const noexcept {
    return *mPtr;
  }

  auto operator<=>(const Iterator &) const = default;

private:
  T *mPtr;
};

template <typename T> struct Span {
  Iterator<const T> begin;
  Iterator<const T> end;
};

template<typename T, size_t NUM_ELEMS>
class SmallVec final {
public:
  template <typename... Args> SmallVec(Args &&...args) noexcept {
    const auto constructSingle = [this]<typename U>(U &&arg) {
      // TODO(javier-varez): add check call to validate size
      std::construct_at(reinterpret_cast<T *>(&mStorage[mUsedSize]),
                        std::forward<U>(arg));
      mUsedSize += 1;
    };

    const auto handleArg = [constructSingle]<typename U>(U &&arg) {
      if constexpr (std::same_as<std::remove_cv_t<std::remove_reference_t<
                                     std::decay_t<U>>>,
                                 Span<const T>>) {
        auto current = arg.begin;
        while (current < arg.end) {
          constructSingle(*current);
          ++current;
        }
      } else {
        constructSingle(std::forward<U>(arg));
      }
    };

    (handleArg(std::forward<Args>(args)), ...);
  }

  SmallVec(const SmallVec& other) {
    for (size_t i = 0; i < other.size(); i++) {
      std::construct_at(reinterpret_cast<T *>(&mStorage[mUsedSize]), other[i]);
      mUsedSize += 1;
    }
  }
  SmallVec(SmallVec& other) {
    for (size_t i = 0; i < other.size(); i++) {
      std::construct_at(reinterpret_cast<T *>(&mStorage[mUsedSize]), other[i]);
      mUsedSize += 1;
    }
  }
  SmallVec(SmallVec&& other) {
    for (size_t i = 0; i < other.size(); i++) {
      std::construct_at(reinterpret_cast<T *>(&mStorage[mUsedSize]), other[i]);
      mUsedSize += 1;
    }
  }

  SmallVec& operator=(const SmallVec& other) {
    if (this != &other) {
      for (size_t i = 0; i < other.size(); i++) {
        std::construct_at(reinterpret_cast<T *>(&mStorage[mUsedSize]), other[i]);
        mUsedSize += 1;
      }
    }
    return *this;
  }

  SmallVec& operator=(SmallVec&& other) {
    if (this != &other) {
      for (size_t i = 0; i < other.size(); i++) {
        std::construct_at(reinterpret_cast<T *>(&mStorage[mUsedSize]), other[i]);
        mUsedSize += 1;
      }
    }
    return *this;
  }

  ~SmallVec() noexcept {
    for (size_t i = 0; i < mUsedSize; ++i) {
      T& elem = *std::launder(reinterpret_cast<T *>(&mStorage[i]));
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
  std::array<Storage, NUM_ELEMS>
      mStorage; // Uninitialized for performance
  size_t mUsedSize{};
};

namespace detail {

template <typename T, typename U, typename... Rest>
static size_t countElems(U &&current, Rest &&...rest) {
  size_t count = 1;
  if constexpr (std::same_as<std::remove_cv_t<
                                 std::remove_reference_t<std::decay_t<U>>>,
                             Span<const T>>) {
    count = current.end - current.begin;
  }

  if constexpr (sizeof...(rest) == 0) {
    return count;
  } else {
    return count + countElems<T>(std::forward<Rest>(rest)...);
  }
}

}  // namespace detail

template<typename T>
class LargeVec final {
public:
  template <typename... Args> LargeVec(Args &&...args) noexcept {
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

private:
  Rc<std::vector<T>> mInner;
};

/**
 * \brief Immutable vector implementation
 */
template <typename T> class Vec {
  constexpr static size_t SMALL_VEC_NUM_ELEMS = 6;

  using SmallVec = SmallVec<T, SMALL_VEC_NUM_ELEMS>;
  using LargeVec = LargeVec<T>;

public:
  Vec() noexcept = default;

  template <typename... Args> static Vec makeVec(Args &&...args) noexcept {
    const size_t count = detail::countElems<T>(std::forward<Args>(args)...);
    if (count > SMALL_VEC_NUM_ELEMS) {
      // Construct new refcounted
      return Vec{std::in_place_type<LargeVec>, std::forward<Args>(args)...};
    }
    return Vec{std::in_place_type<SmallVec>, std::forward<Args>(args)...};
  }

  Vec(const Vec &) = default;
  Vec(Vec &&) = default;
  Vec &operator=(const Vec &) = default;
  Vec &operator=(Vec &&) = default;

  constexpr Iterator<const T> begin() const noexcept {
    return std::visit(
        []<typename I>(const I &elem) {
          return elem.begin();
        },
        mInner);
  }

  constexpr Iterator<const T> end() const noexcept {
    return std::visit(
        []<typename I>(const I &elem) {
          return elem.end();
        },
        mInner);
  }

  constexpr const T &operator[](size_t index) const noexcept {
    return std::visit(
        [index]<typename I>(const I &elem) -> const T & {
          return elem[index];
        },
        mInner);
  }

  constexpr size_t size() const noexcept {
    return std::visit(
        []<typename I>(const I &elem) {
          return elem.size();
        },
        mInner);
  }

  constexpr size_t capacity() const noexcept {
    return std::visit(
        []<typename I>(const I &elem) {
          return elem.capacity();
        },
        mInner);
  }

  template <typename... Args>
  constexpr Vec copyAppend(Args &&...args) const noexcept {
    // Decide what type of object to create
    const size_t count = detail::countElems<T>(std::forward<Args>(args)...) + size();
    const Span<const T> currElemsSpan{.begin = begin(), .end = end()};

    if (count > SMALL_VEC_NUM_ELEMS) {
      // Construct new refcounted
      return Vec{std::in_place_type<LargeVec>, currElemsSpan,
                 std::forward<Args>(args)...};
    }

    return Vec{std::in_place_type<SmallVec>, currElemsSpan,
               std::forward<Args>(args)...};
  }

  bool isSmallVec() const noexcept { return mInner.index() == 0; }

private:
  std::variant<SmallVec, LargeVec> mInner;

  template <typename... Args>
  Vec(std::in_place_type_t<SmallVec>, Args &&...args) noexcept
      : mInner(std::in_place_type<SmallVec>, std::forward<Args>(args)...) { }

  template <typename... Args>
  Vec(std::in_place_type_t<LargeVec>,
      Args &&...args) noexcept
      : mInner(std::in_place_type<LargeVec>, std::forward<Args>(args)...) { }
};

} // namespace runtime

#endif // _RUNTIME_VECTOR_H_
