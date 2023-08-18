#include <gtest/gtest.h>

#include <box.h>

namespace {

struct Stats {
  size_t constructed{};
  size_t copyConstructed{};
  size_t moveConstructed{};
  size_t copyAssigned{};
  size_t moveAssigned{};
  size_t destructed{};
};

class StatCounter {
public:
  StatCounter(Stats &stats) noexcept : mStats{stats} { ++mStats.constructed; }
  StatCounter(const StatCounter &other) noexcept : mStats(other.mStats) {
    ++mStats.copyConstructed;
  }
  StatCounter(StatCounter &&other) noexcept : mStats(other.mStats) {
    ++mStats.moveConstructed;
  }
  StatCounter &operator=(const StatCounter &) noexcept {
    ++mStats.copyAssigned;
    return *this;
  }
  StatCounter &operator=(StatCounter &&) noexcept {
    ++mStats.moveAssigned;
    return *this;
  }
  ~StatCounter() noexcept { ++mStats.destructed; }

private:
  Stats &mStats;
};

class Base {
public:
  virtual int doStuff() noexcept = 0;
  virtual ~Base() noexcept = default;
};

class Derived final : public Base {
public:
  int doStuff() noexcept final { return 123; }
  ~Derived() noexcept final = default;
};

}  // namespace

TEST(Box, ConstructionAndDestruction) {
  Stats stats{};
  EXPECT_EQ(stats.constructed, 0);
  EXPECT_EQ(stats.copyConstructed, 0);
  EXPECT_EQ(stats.moveConstructed, 0);
  EXPECT_EQ(stats.copyAssigned, 0);
  EXPECT_EQ(stats.moveAssigned, 0);
  EXPECT_EQ(stats.destructed, 0);

  {
    runtime::Box<StatCounter> box{runtime::Marker<StatCounter>{}, stats};

    EXPECT_EQ(stats.constructed, 1);
    EXPECT_EQ(stats.copyConstructed, 0);
    EXPECT_EQ(stats.moveConstructed, 0);
    EXPECT_EQ(stats.copyAssigned, 0);
    EXPECT_EQ(stats.moveAssigned, 0);
    EXPECT_EQ(stats.destructed, 0);
  }

  EXPECT_EQ(stats.constructed, 1);
  EXPECT_EQ(stats.copyConstructed, 0);
  EXPECT_EQ(stats.moveConstructed, 0);
  EXPECT_EQ(stats.copyAssigned, 0);
  EXPECT_EQ(stats.moveAssigned, 0);
  EXPECT_EQ(stats.destructed, 1);
}

TEST(Box, AccessInner) {
  runtime::Box<int> box{runtime::Marker<int>{}, 10};
  EXPECT_EQ(*box.operator->(), 10);
  EXPECT_EQ(box.operator*(), 10);
  EXPECT_EQ(*box, 10);

  *box = 12;

  EXPECT_EQ(*box.operator->(), 12);
  EXPECT_EQ(box.operator*(), 12);
  EXPECT_EQ(*box, 12);

  const runtime::Box<int> constBox{runtime::Marker<int>{}, 12};
  EXPECT_EQ(*constBox.operator->(), 12);
  EXPECT_EQ(constBox.operator*(), 12);
}

TEST(Box, FromDerived) {
  runtime::Box<Base> box{runtime::Marker<Derived>{}};
  EXPECT_EQ(box->doStuff(), 123);
}
