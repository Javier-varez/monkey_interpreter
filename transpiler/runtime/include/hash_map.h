#pragma once

#include <object.h>
#include <rc.h>

#include <functional>

namespace runtime {

class HashMap final {
 public:
  struct KvPair {
    const Object& k;
    const Object& v;
  };

  HashMap() noexcept;

  template <typename... Args>
    requires((std::same_as<Args, KvPair> && ...))
  HashMap(const Args&... args) noexcept : HashMap() {
    (pushKvPair(args), ...);
  }

  const Object& operator[](const Object& key) const noexcept;

  void forEach(const std::function<void(const Object&, const Object&)>&
                   callable) const noexcept;

  ~HashMap() noexcept;

 private:
  class Impl;
  Rc<Impl> mImpl;

  void pushKvPair(const KvPair& pair) noexcept;
};

}  // namespace runtime
