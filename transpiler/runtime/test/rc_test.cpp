#include <gtest/gtest.h>

#include <rc.h>

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

TEST(Rc, ConstructorAndCopyConstruction) {
  Stats stats{};
  EXPECT_EQ(stats.constructed, 0);
  EXPECT_EQ(stats.copyConstructed, 0);
  EXPECT_EQ(stats.moveConstructed, 0);
  EXPECT_EQ(stats.copyAssigned, 0);
  EXPECT_EQ(stats.moveAssigned, 0);
  EXPECT_EQ(stats.destructed, 0);

  {
    runtime::Rc<StatCounter> rc{stats};

    EXPECT_EQ(stats.constructed, 1);
    EXPECT_EQ(stats.copyConstructed, 0);
    EXPECT_EQ(stats.moveConstructed, 0);
    EXPECT_EQ(stats.copyAssigned, 0);
    EXPECT_EQ(stats.moveAssigned, 0);
    EXPECT_EQ(stats.destructed, 0);

    {
      runtime::Rc<StatCounter> copy{rc};

      EXPECT_EQ(stats.constructed, 1);
      EXPECT_EQ(stats.copyConstructed, 0);
      EXPECT_EQ(stats.moveConstructed, 0);
      EXPECT_EQ(stats.copyAssigned, 0);
      EXPECT_EQ(stats.moveAssigned, 0);
      EXPECT_EQ(stats.destructed, 0);
      EXPECT_EQ(rc.operator->(), copy.operator->());
      EXPECT_EQ(&rc.operator*(), &copy.operator*());

      {
        const runtime::Rc<StatCounter> &ref{rc};
        runtime::Rc<StatCounter> anotherCopy{ref};

        EXPECT_EQ(stats.constructed, 1);
        EXPECT_EQ(stats.copyConstructed, 0);
        EXPECT_EQ(stats.moveConstructed, 0);
        EXPECT_EQ(stats.copyAssigned, 0);
        EXPECT_EQ(stats.moveAssigned, 0);
        EXPECT_EQ(stats.destructed, 0);

        EXPECT_EQ(rc.operator->(), anotherCopy.operator->());
        EXPECT_EQ(&rc.operator*(), &anotherCopy.operator*());
      }
    }

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

TEST(Rc, AssignmentOp) {
  Stats stats{};
  EXPECT_EQ(stats.constructed, 0);
  EXPECT_EQ(stats.copyConstructed, 0);
  EXPECT_EQ(stats.moveConstructed, 0);
  EXPECT_EQ(stats.copyAssigned, 0);
  EXPECT_EQ(stats.moveAssigned, 0);
  EXPECT_EQ(stats.destructed, 0);

  Stats innerStats{};

  {
    runtime::Rc<StatCounter> rc{stats};

    EXPECT_EQ(stats.constructed, 1);
    EXPECT_EQ(stats.copyConstructed, 0);
    EXPECT_EQ(stats.moveConstructed, 0);
    EXPECT_EQ(stats.copyAssigned, 0);
    EXPECT_EQ(stats.moveAssigned, 0);
    EXPECT_EQ(stats.destructed, 0);

    {
      runtime::Rc<StatCounter> rc2{innerStats};

      EXPECT_EQ(innerStats.constructed, 1);
      EXPECT_EQ(innerStats.copyConstructed, 0);
      EXPECT_EQ(innerStats.moveConstructed, 0);
      EXPECT_EQ(innerStats.copyAssigned, 0);
      EXPECT_EQ(innerStats.moveAssigned, 0);
      EXPECT_EQ(innerStats.destructed, 0);
      EXPECT_NE(rc.operator->(), rc2.operator->());
      EXPECT_NE(&rc.operator*(), &rc2.operator*());

      EXPECT_EQ(stats.constructed, 1);
      EXPECT_EQ(stats.copyConstructed, 0);
      EXPECT_EQ(stats.moveConstructed, 0);
      EXPECT_EQ(stats.copyAssigned, 0);
      EXPECT_EQ(stats.moveAssigned, 0);
      EXPECT_EQ(stats.destructed, 0);

      // Now copy assign one to the other
      rc2 = rc;

      EXPECT_EQ(innerStats.constructed, 1);
      EXPECT_EQ(innerStats.copyConstructed, 0);
      EXPECT_EQ(innerStats.moveConstructed, 0);
      EXPECT_EQ(innerStats.copyAssigned, 0);
      EXPECT_EQ(innerStats.moveAssigned, 0);
      EXPECT_EQ(innerStats.destructed, 1);
      EXPECT_EQ(rc.operator->(), rc2.operator->());
      EXPECT_EQ(&rc.operator*(), &rc2.operator*());
      EXPECT_EQ(stats.constructed, 1);
      EXPECT_EQ(stats.copyConstructed, 0);
      EXPECT_EQ(stats.moveConstructed, 0);
      EXPECT_EQ(stats.copyAssigned, 0);
      EXPECT_EQ(stats.moveAssigned, 0);
      EXPECT_EQ(stats.destructed, 0);
    }

    EXPECT_EQ(innerStats.constructed, 1);
    EXPECT_EQ(innerStats.copyConstructed, 0);
    EXPECT_EQ(innerStats.moveConstructed, 0);
    EXPECT_EQ(innerStats.copyAssigned, 0);
    EXPECT_EQ(innerStats.moveAssigned, 0);
    EXPECT_EQ(innerStats.destructed, 1);
    EXPECT_EQ(stats.constructed, 1);
    EXPECT_EQ(stats.copyConstructed, 0);
    EXPECT_EQ(stats.moveConstructed, 0);
    EXPECT_EQ(stats.copyAssigned, 0);
    EXPECT_EQ(stats.moveAssigned, 0);
    EXPECT_EQ(stats.destructed, 0);
  }

  EXPECT_EQ(innerStats.constructed, 1);
  EXPECT_EQ(innerStats.copyConstructed, 0);
  EXPECT_EQ(innerStats.moveConstructed, 0);
  EXPECT_EQ(innerStats.copyAssigned, 0);
  EXPECT_EQ(innerStats.moveAssigned, 0);
  EXPECT_EQ(innerStats.destructed, 1);
  EXPECT_EQ(stats.constructed, 1);
  EXPECT_EQ(stats.copyConstructed, 0);
  EXPECT_EQ(stats.moveConstructed, 0);
  EXPECT_EQ(stats.copyAssigned, 0);
  EXPECT_EQ(stats.moveAssigned, 0);
  EXPECT_EQ(stats.destructed, 1);
}
