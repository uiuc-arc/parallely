import unittest
import rely_recover
from rely_recover import Constraint


def sameConstraint(c1, c2):
    if len(c1) != len(c2):
        return False
    for i in range(len(c1)):
        same_l = (c1[i].limit == c2[i].limit)
        same_c = (c1[i].condition == c2[i].condition)
        print "*************************", c1[i].multiplicative, c2[i].multiplicative, c1[i].multiplicative - c2[i].multiplicative
        same_mult = abs(c1[i].multiplicative - c2[i].multiplicative) <= 0.00001
        same_set = (c1[i].jointreliability == c2[i].jointreliability)
        print same_l, same_c, same_mult, same_set
        if not (same_l and same_c and same_mult and same_set):
            return False
    return True


class TestSimple(unittest.TestCase):

    def test_basic(self):
        programfile = open("tests/basic/basic1.code", 'r')
        spec = open("tests/basic/basic1.spec", 'r').read()
        unroll = 0
        program_str = programfile.read()
        expected_result = [Constraint(limit=0.99,
                                      condition='<=',
                                      multiplicative=0.9,
                                      jointreliability=set([]))]
        new_result = rely_recover.main(program_str, spec, unroll)
        assert(sameConstraint(new_result, expected_result))

    def test_basic2(self):
        programfile = open("tests/basic/basic2.code", 'r')
        spec = open("tests/basic/basic1.spec", 'r').read()
        unroll = 0
        program_str = programfile.read()
        expected_result = [Constraint(limit=0.99,
                                      condition='<=',
                                      multiplicative=0.85,
                                      jointreliability=set([]))]
        new_result = rely_recover.main(program_str, spec, unroll)
        assert(sameConstraint(new_result, expected_result))

    def test_basic3(self):
        programfile = open("tests/basic/basic3.code", 'r')
        spec = open("tests/basic/basic1.spec", 'r').read()
        unroll = 0
        program_str = programfile.read()
        expected_result = [Constraint(limit=0.99,
                                      condition='<=',
                                      multiplicative=0.09,
                                      jointreliability=set(['temp2', 'outNFloat']))]
        new_result = rely_recover.main(program_str, spec, unroll)
        assert(sameConstraint(new_result, expected_result))

    def test_basic4(self):
        programfile = open("tests/basic/basic4.code", 'r')
        spec = open("tests/basic/basic1.spec", 'r').read()
        unroll = 0
        program_str = programfile.read()
        expected_result = [Constraint(limit=0.99,
                                      condition='<=',
                                      multiplicative=0.36125,
                                      jointreliability=set(['outNFloat']))]
        new_result = rely_recover.main(program_str, spec, unroll)
        assert(sameConstraint(new_result, expected_result))


class TestRecovery(unittest.TestCase):
    def testRecovery1(self):
        programfile = open("tests/recovery/recovery1.code", 'r')
        spec = open("tests/recovery/recovery.spec", 'r').read()
        unroll = 0
        program_str = programfile.read()
        expected_result = [Constraint(limit=0.99,
                                      condition='<=',
                                      multiplicative=0.99,
                                      jointreliability=set([]))]
        new_result = rely_recover.main(program_str, spec, unroll)
        assert(sameConstraint(new_result, expected_result))

    def testRecovery2(self):
        programfile = open("tests/recovery/recovery2.code", 'r')
        spec = open("tests/recovery/recovery.spec", 'r').read()
        unroll = 0
        program_str = programfile.read()
        expected_result = [Constraint(limit=0.99,
                                      condition='<=',
                                      multiplicative=0.945,
                                      jointreliability=set([]))]
        new_result = rely_recover.main(program_str, spec, unroll)
        assert(sameConstraint(new_result, expected_result))


if __name__ == '__main__':
    unittest.main()
