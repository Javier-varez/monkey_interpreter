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

/**
 * \brief Immutable vector implementation
 */
template <typename T> class Vec {
  constexpr static size_t SMALL_VEC_NUM_ELEMS = 6;

public:
  Vec() noexcept = default;

  template <typename... Args> static Vec makeVec(Args &&...args) noexcept {
    // Decide what type of object to create
    const size_t count = countElems(std::forward<Args>(args)...);

    if (count > SMALL_VEC_NUM_ELEMS) {
      // Construct new refcounted
      return Vec{std::in_place_type<LargeVec>, count,
                 std::forward<Args>(args)...};
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
          if constexpr (std::same_as<I, SmallVec>) {
            return elem.begin();
          } else {
            return Iterator{&*elem->begin()};
          }
        },
        mInner);
  }

  constexpr Iterator<const T> end() const noexcept {
    return std::visit(
        []<typename I>(const I &elem) {
          if constexpr (std::same_as<I, SmallVec>) {
            return elem.end();
          } else {
            return Iterator{&*elem->end()};
          }
        },
        mInner);
  }

  constexpr const T &operator[](size_t index) const noexcept {
    return std::visit(
        [index]<typename I>(const I &elem) -> const T & {
          if constexpr (std::same_as<I, SmallVec>) {
            return elem[index];
          } else {
            return (*elem)[index];
          }
        },
        mInner);
  }

  constexpr size_t size() const noexcept {
    return std::visit(
        []<typename I>(const I &elem) {
          if constexpr (std::same_as<I, SmallVec>) {
            return elem.size();
          } else {
            return elem->size();
          }
        },
        mInner);
  }

  constexpr size_t capacity() const noexcept {
    return std::visit(
        []<typename I>(const I &elem) {
          if constexpr (std::same_as<I, SmallVec>) {
            return elem.capacity();
          } else {
            return elem->capacity();
          }
        },
        mInner);
  }

  template <typename... Args>
  constexpr Vec copyAppend(Args &&...args) const noexcept {
    // Decide what type of object to create
    const size_t count = countElems(std::forward<Args>(args)...) + size();
    const Span<const T> currElemsSpan{.begin = begin(), .end = end()};

    if (count > SMALL_VEC_NUM_ELEMS) {
      // Construct new refcounted
      return Vec{std::in_place_type<LargeVec>, count, currElemsSpan,
                 std::forward<Args>(args)...};
    }

    return Vec{std::in_place_type<SmallVec>, currElemsSpan,
               std::forward<Args>(args)...};
  }

  bool isSmallVec() const noexcept { return mInner.index() == 0; }

private:
  class SmallVec {
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

    constexpr size_t capacity() const noexcept { return SMALL_VEC_NUM_ELEMS; }

    constexpr size_t size() const noexcept { return mUsedSize; }

  private:
    using Storage = std::aligned_storage_t<sizeof(T), alignof(T)>;
    std::array<Storage, SMALL_VEC_NUM_ELEMS>
        mStorage; // Uninitialized for performance
    size_t mUsedSize{};
  };

  using LargeVec = Rc<std::vector<T>>;

  std::variant<SmallVec, LargeVec> mInner;

  template <typename U, typename... Rest>
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
      return count + countElems(std::forward<Rest>(rest)...);
    }
  }

  template <typename... Args>
  Vec(std::in_place_type_t<SmallVec>, Args &&...args) noexcept
      : mInner(std::in_place_type<SmallVec>, std::forward<Args>(args)...) {}

  template <typename... Args>
  Vec(std::in_place_type_t<LargeVec>, const size_t numElems,
      Args &&...args) noexcept
      : mInner(std::in_place_type<LargeVec>) {
    LargeVec &vec = std::get<LargeVec>(mInner);
    vec->reserve(numElems);

    const auto handleArg = [&vec]<typename U>(U &&arg) {
      if constexpr (std::same_as<std::remove_cv_t<
                                     std::remove_reference_t<std::decay_t<U>>>,
                                 Span<const T>>) {
        Iterator<const T> current = arg.begin;
        while (current < arg.end) {
          vec->push_back(*current);
          ++current;
        }
      } else {
        vec->push_back(std::forward<U>(arg));
      }
    };

    (handleArg(std::forward<Args>(args)), ...);
  }
};

} // namespace runtime

#endif // _RUNTIME_VECTOR_H_
