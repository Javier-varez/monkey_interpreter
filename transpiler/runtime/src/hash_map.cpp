#include <hash_map.h>
#include <object.h>
#include <var_args.h>

#include <functional>
#include <unordered_map>

namespace runtime {

namespace {

struct ObjectWrapper {
  const Object& obj;
};

[[nodiscard]] bool operator==(const ObjectWrapper& lhs,
                              const ObjectWrapper& rhs) noexcept {
  return lhs.obj.equals(rhs.obj);
}

[[nodiscard]] bool operator!=(const ObjectWrapper& lhs,
                              const ObjectWrapper& rhs) noexcept {
  return !lhs.obj.equals(rhs.obj);
}

[[nodiscard]] bool operator==(const ObjectWrapper& lhs,
                              const Object& rhs) noexcept {
  return lhs.obj.equals(rhs);
}

[[nodiscard]] bool operator==(const Object& lhs,
                              const ObjectWrapper& rhs) noexcept {
  return lhs.equals(rhs.obj);
}

[[nodiscard]] bool operator!=(const ObjectWrapper& lhs,
                              const Object& rhs) noexcept {
  return !lhs.obj.equals(rhs);
}

[[nodiscard]] bool operator!=(const Object& lhs,
                              const ObjectWrapper& rhs) noexcept {
  return !lhs.equals(rhs.obj);
}

}  // namespace

}  // namespace runtime

namespace std {

template <>
struct hash<runtime::ObjectWrapper> final {
  std::int64_t operator()(
      const runtime::ObjectWrapper& wrapper) const noexcept {
    return wrapper.obj.hash();
  }
};

}  // namespace std

namespace runtime {

class HashMap::Impl {
 public:
  Impl() noexcept;

  void pushKvPair(const KvPair& pair) noexcept;

  const Object& operator[](const Object& key) const noexcept;

  void forEach(const std::function<void(const Object&, const Object&)>&
                   callable) const noexcept;

 private:
  std::unordered_map<ObjectWrapper, Object> mMap;
};

HashMap::Impl::Impl() noexcept {}

void HashMap::Impl::pushKvPair(const KvPair& pair) noexcept {
  mMap[ObjectWrapper{pair.k}] = pair.v;
}

const Object& HashMap::Impl::operator[](const Object& key) const noexcept {
  const ObjectWrapper k{key};
  if (mMap.contains(k)) {
    return mMap.at(k);
  }
  return Object::nil();
}

void HashMap::Impl::forEach(
    const std::function<void(const Object&, const Object&)>& callable)
    const noexcept {
  for (const auto& [k, v] : mMap) {
    callable(k.obj, v);
  }
}

HashMap::HashMap() noexcept : mImpl{} {}

const Object& HashMap::operator[](const Object& key) const noexcept {
  return (*mImpl)[key];
}

HashMap::~HashMap() noexcept {}

void HashMap::pushKvPair(const KvPair& pair) noexcept {
  mImpl->pushKvPair(pair);
}

void HashMap::forEach(const std::function<void(const Object&, const Object&)>&
                          callable) const noexcept {
  mImpl->forEach(callable);
}

}  // namespace runtime
