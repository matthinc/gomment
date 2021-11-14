import unittest

from idx_hierarchy import idx_to_hierarchy, hierarchy_to_idx, hierarchy_to_touched_at

class TestHierarchy(unittest.TestCase):

    def test_first_idx(self):
        dim = [4, 3, 2]

        self.assertEqual([0], idx_to_hierarchy(dim, 0))
        self.assertEqual([1], idx_to_hierarchy(dim, 1))
        self.assertEqual([2], idx_to_hierarchy(dim, 2))
        self.assertEqual([3], idx_to_hierarchy(dim, 3))

    def test_second_idx(self):
        dim = [4, 3, 2]
        offset = dim[0]

        for x in range(4):
            for y in range(3):
                self.assertEqual([x, y], idx_to_hierarchy(dim, offset + 3 * x + y))

    def test_third_idx(self):
        dim = [4, 3, 2]
        offset = dim[0] + (dim[0] * dim[1])

        for x in range(4):
            for y in range(3):
                for z in range(2):
                    self.assertEqual([x, y, z], idx_to_hierarchy(dim, offset + 6 * x + 2 * y + z))


    def test_first_hie(self):
        dim = [4, 3, 2]

        self.assertEqual(0, hierarchy_to_idx(dim, [0]))
        self.assertEqual(1, hierarchy_to_idx(dim, [1]))
        self.assertEqual(2, hierarchy_to_idx(dim, [2]))
        self.assertEqual(3, hierarchy_to_idx(dim, [3]))

    def test_second_hie(self):
        dim = [4, 3, 2]
        offset = dim[0]

        # [3, 1]
        for x in range(4):
            for y in range(3):
                self.assertEqual(offset + 3 * x + y, hierarchy_to_idx(dim, [x, y]))


    def test_third_hie(self):
        dim = [4, 3, 2]
        offset = dim[0] + (dim[0] * dim[1])

        # [6, 2, 1]
        for x in range(4):
            for y in range(3):
                for z in range(2):
                    self.assertEqual(offset + 6 * x + 2 * y + z, hierarchy_to_idx(dim, [x, y, z]))

    def test_touched_at_2_2(self):
        dim = [2, 2]

        # [0] = 3
        # [0, 0] = 2
        # [0, 1] = 3
        # [1] = 5
        # [1, 0] = 4
        # [1, 1] = 5

        self.assertEqual(2, hierarchy_to_touched_at(dim, [0, 0]))
        self.assertEqual(3, hierarchy_to_touched_at(dim, [0, 1]))
        self.assertEqual(4, hierarchy_to_touched_at(dim, [1, 0]))
        self.assertEqual(5, hierarchy_to_touched_at(dim, [1, 1]))
        self.assertEqual(3, hierarchy_to_touched_at(dim, [0]))
        self.assertEqual(5, hierarchy_to_touched_at(dim, [1]))

    def test_touched_at_2_2_2(self):
        dim = [2, 2, 2]

        # [0] = 9
        # [0, 0] = 7
        # [0, 0, 0] = 6
        # [0, 0, 1] = 7
        # [0, 1] = 9
        # [0, 1, 0] = 8
        # [0, 1, 1] = 9
        # [1] = 13
        # [1, 0] = 11
        # [1, 0, 0] = 10
        # [1, 0, 1] = 11
        # [1, 1] = 13
        # [1, 1, 0] = 12
        # [1, 1, 1] = 13

        self.assertEqual(6 , hierarchy_to_touched_at(dim, [0, 0, 0]))
        self.assertEqual(7 , hierarchy_to_touched_at(dim, [0, 0, 1]))
        self.assertEqual(8 , hierarchy_to_touched_at(dim, [0, 1, 0]))
        self.assertEqual(9 , hierarchy_to_touched_at(dim, [0, 1, 1]))
        self.assertEqual(10, hierarchy_to_touched_at(dim, [1, 0, 0]))
        self.assertEqual(11, hierarchy_to_touched_at(dim, [1, 0, 1]))
        self.assertEqual(12, hierarchy_to_touched_at(dim, [1, 1, 0]))
        self.assertEqual(13, hierarchy_to_touched_at(dim, [1, 1, 1]))
        self.assertEqual(7 , hierarchy_to_touched_at(dim, [0, 0]))
        self.assertEqual(9 , hierarchy_to_touched_at(dim, [0, 1]))
        self.assertEqual(11, hierarchy_to_touched_at(dim, [1, 0]))
        self.assertEqual(13, hierarchy_to_touched_at(dim, [1, 1]))
        self.assertEqual(9 , hierarchy_to_touched_at(dim, [0]))
        self.assertEqual(13, hierarchy_to_touched_at(dim, [1]))


if __name__ == '__main__':
    unittest.main()
