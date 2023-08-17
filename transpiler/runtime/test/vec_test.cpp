#include <gtest/gtest.h>

#include <vec.h>

namespace {

template<typename T>
T makeRange(const size_t start, const size_t end) {
  const auto abs = [](auto arg) {
    if (arg < 0)
      return -arg;
    return arg;
  };
  const size_t sizeHint = abs(end - start);

  int64_t current = start;

  if (start > end) {
    return T {[&current, end]() -> std::optional<int> {
      if (current > end) {
        return std::optional<int>{current--};
      }
      return {};
    }, sizeHint};
  }

  return T {[&current, end]() -> std::optional<int> {
    if (current < end) {
      return std::optional<int>{current++};
    }
    return {};
  }, sizeHint};
}

}  // namespace

TEST(Vec, ConstructEmpty) {
  const runtime::Vec<int> v;

  EXPECT_EQ(v.size(), 0);
  EXPECT_TRUE(v.isSmallVec());
}

TEST(Vec, MakeVec) {
  {
    const auto v = runtime::Vec<int>::makeVec(1, 2, 3, 4, 5, 6);
    EXPECT_EQ(v.size(), 6);
    for (int i = 0; i < v.size(); i++) {
      EXPECT_EQ(v[i], i + 1);
    }
    EXPECT_TRUE(v.isSmallVec());
  }

  {
    const auto v = runtime::Vec<int>::makeVec(1, 2, 3, 4, 5, 6, 7);
    EXPECT_EQ(v.size(), 7);
    for (int i = 0; i < v.size(); i++) {
      EXPECT_EQ(v[i], i + 1);
    }
    EXPECT_FALSE(v.isSmallVec());
  }
}

TEST(Vec, Iterators) {
  {
    const auto v = runtime::Vec<int>::makeVec(1, 2, 3, 4, 5, 6);
    EXPECT_EQ(v.size(), 6);
    int expected = 1;
    for (const int &val : v) {
      EXPECT_EQ(val, expected);
      expected++;
    }
    EXPECT_EQ(expected, 7);
    EXPECT_TRUE(v.isSmallVec());
  }

  {
    const auto v = runtime::Vec<int>::makeVec(1, 2, 3, 4, 5, 6, 7);
    EXPECT_EQ(v.size(), 7);
    int expected = 1;
    for (const int &val : v) {
      EXPECT_EQ(val, expected);
      expected++;
    }
    EXPECT_EQ(expected, 8);
    EXPECT_FALSE(v.isSmallVec());
  }
}

TEST(Vec, CopyAppend) {
  const runtime::Vec<int> v;

  EXPECT_EQ(v.size(), 0);
  EXPECT_TRUE(v.isSmallVec());

  {
    const runtime::Vec<int> v2{v.copyAppend(1, 2, 3)};
    EXPECT_EQ(v2.size(), 3);
    EXPECT_EQ(v2[0], 1);
    EXPECT_EQ(v2[1], 2);
    EXPECT_EQ(v2[2], 3);
    EXPECT_TRUE(v2.isSmallVec());
  }

  {
    const runtime::Vec<int> v2{v.copyAppend(1, 2, 3).copyAppend(1, 2, 3)};
    EXPECT_EQ(v2.size(), 6);
    EXPECT_EQ(v2[0], 1);
    EXPECT_EQ(v2[1], 2);
    EXPECT_EQ(v2[2], 3);
    EXPECT_EQ(v2[3], 1);
    EXPECT_EQ(v2[4], 2);
    EXPECT_EQ(v2[5], 3);
    EXPECT_TRUE(v2.isSmallVec());
  }

  {
    const runtime::Vec<int> v2{v.copyAppend(1, 2, 3).copyAppend(1, 2, 3, 4)};
    EXPECT_EQ(v2.size(), 7);
    EXPECT_EQ(v2[0], 1);
    EXPECT_EQ(v2[1], 2);
    EXPECT_EQ(v2[2], 3);
    EXPECT_EQ(v2[3], 1);
    EXPECT_EQ(v2[4], 2);
    EXPECT_EQ(v2[5], 3);
    EXPECT_EQ(v2[6], 4);
    EXPECT_FALSE(v2.isSmallVec());
  }
}

TEST(Vec, Copy) {
  {
    const auto v = runtime::Vec<int>::makeVec(1, 2, 3, 4, 5);

    EXPECT_EQ(v.size(), 5);
    EXPECT_TRUE(v.isSmallVec());

    {
      const runtime::Vec<int> v2{v};
      EXPECT_EQ(v2.size(), 5);
      EXPECT_EQ(v2[0], 1);
      EXPECT_EQ(v2[1], 2);
      EXPECT_EQ(v2[2], 3);
      EXPECT_EQ(v2[3], 4);
      EXPECT_EQ(v2[4], 5);

      // Elements are allocated statically, copied over on copy
      EXPECT_NE(&v[0], &v2[0]);
      EXPECT_TRUE(v2.isSmallVec());
    }
  }

  {
    const auto v = runtime::Vec<int>::makeVec(1, 2, 3, 4, 5, 6, 7);

    EXPECT_EQ(v.size(), 7);
    EXPECT_FALSE(v.isSmallVec());

    {
      const runtime::Vec<int> v2{v};
      EXPECT_EQ(v2.size(), 7);
      EXPECT_EQ(v2[0], 1);
      EXPECT_EQ(v2[1], 2);
      EXPECT_EQ(v2[2], 3);
      EXPECT_EQ(v2[3], 4);
      EXPECT_EQ(v2[4], 5);
      EXPECT_EQ(v2[5], 6);
      EXPECT_EQ(v2[6], 7);

      // Elements are allocated dynamically, reference counting the actual
      // immutable vector
      EXPECT_EQ(&v[0], &v2[0]);
      EXPECT_FALSE(v2.isSmallVec());
    }
  }
}

TEST(Vec, CopyAssign) {
  {
    const auto v = runtime::Vec<int>::makeVec(1, 2, 3, 4, 5);

    EXPECT_EQ(v.size(), 5);
    EXPECT_TRUE(v.isSmallVec());

    {
      runtime::Vec<int> v2;
      v2 = v;

      EXPECT_EQ(v2.size(), 5);
      EXPECT_EQ(v2[0], 1);
      EXPECT_EQ(v2[1], 2);
      EXPECT_EQ(v2[2], 3);
      EXPECT_EQ(v2[3], 4);
      EXPECT_EQ(v2[4], 5);

      // Elements are allocated statically, copied over on copy
      EXPECT_NE(&v[0], &v2[0]);
      EXPECT_TRUE(v2.isSmallVec());
    }
  }

  {
    const auto v = runtime::Vec<int>::makeVec(1, 2, 3, 4, 5, 6, 7);

    EXPECT_EQ(v.size(), 7);
    EXPECT_FALSE(v.isSmallVec());

    {
      runtime::Vec<int> v2;
      v2 = v;

      EXPECT_EQ(v2.size(), 7);
      EXPECT_EQ(v2[0], 1);
      EXPECT_EQ(v2[1], 2);
      EXPECT_EQ(v2[2], 3);
      EXPECT_EQ(v2[3], 4);
      EXPECT_EQ(v2[4], 5);
      EXPECT_EQ(v2[5], 6);
      EXPECT_EQ(v2[6], 7);

      // Elements are allocated dynamically, reference counting the actual
      // immutable vector
      EXPECT_EQ(&v[0], &v2[0]);
      EXPECT_FALSE(v2.isSmallVec());
    }
  }
}

TEST(Vec, MakeWithCallable) {
  {
    runtime::LargeVec vec = makeRange<runtime::LargeVec<int>>(0, 100);
    EXPECT_EQ(vec.size(), 100);
    for (size_t i = 0; i < vec.size(); i++) {
      EXPECT_EQ(vec[i], i);
    }
  }

  {
    runtime::SmallVec vec = makeRange<runtime::SmallVec<int, 100>>(0, 100);
    EXPECT_EQ(vec.size(), 100);
    for (size_t i = 0; i < vec.size(); i++) {
      EXPECT_EQ(vec[i], i);
    }
  }

  {
    runtime::Vec vec = makeRange<runtime::Vec<int>>(0, 6);
    EXPECT_TRUE(vec.isSmallVec());
    EXPECT_EQ(vec.size(), 6);
    for (size_t i = 0; i < vec.size(); i++) {
      EXPECT_EQ(vec[i], i);
    }
  }

  {
    runtime::Vec vec = makeRange<runtime::Vec<int>>(0, 7);
    EXPECT_FALSE(vec.isSmallVec());
    EXPECT_EQ(vec.size(), 7);
    for (size_t i = 0; i < vec.size(); i++) {
      EXPECT_EQ(vec[i], i);
    }
  }
}
