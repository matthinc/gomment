import unittest
import os
import concurrent.futures
import argparse
import pickle
import subprocess


SCRIPT_DIR = os.path.abspath(os.path.dirname(os.path.realpath(__file__)))


def test_suite_to_fq_name(test_suite):
    return f'{test_suite.__class__.__module__}.{test_suite.__class__.__qualname__}.{test_suite._testMethodName}'


class TestRun():
    def __init__(self, fq_test_name):
        self.fq_test_name = fq_test_name

    def run(self):
        loader = unittest.TestLoader()
        suite = loader.loadTestsFromName(self.fq_test_name)

        runner = unittest.TextTestRunner()
        result = runner.run(suite)


class TestRunDispatcher():
    def __init__(self, test_suite):
        self.test_suite = test_suite

    def __str__(self):
        return test_suite_to_fq_name(self.test_suite)

    def run(self):
        state_dir = f'{SCRIPT_DIR}/state/{self}'

        if not os.path.exists(state_dir):
            os.makedirs(state_dir)

        docker_cmd = f'''
        docker run
          --rm
          -u$(id -u):$(id -g)
          -v{SCRIPT_DIR}:/app/test
          -v{state_dir}:/app/test-state
          -w /app
          --env DB_PATH=/app/test-state/test.db
          gomment-test python3 test/system_test_runner.py --direct --tests {self}
        '''.strip().replace('\n', ' ')

        pipe = subprocess.Popen(
            ['/bin/bash', '-c', docker_cmd],
            # stdout=subprocess.PIPE,
            # stderr=subprocess.PIPE,
        )
        output = pipe.communicate()[0]
        return output


class SystemTestRunner():
    def __init__(self, direct, test_names, original_cwd):
        self.direct = direct
        self.test_names = test_names
        self.test_suites = []
        self.original_cwd = original_cwd

    # convert hierarchical test suites to flat test methods
    # https://stackoverflow.com/a/16823380
    def test_suites_to_methods(self, s):
      for test in s:
        if unittest.suite._isnotsuite(test):
          yield test
        else:
          for t in self.test_suites_to_methods(test):
            yield t

    def get_all_test_suites(self):
        os.chdir(SCRIPT_DIR)

        # recursively get all test cases
        loader = unittest.TestLoader()
        suite = loader.discover('.', pattern='test_*.py')

        return list(self.test_suites_to_methods(suite))

    def test_names_to_test_suites(self, test_names):
        os.chdir(SCRIPT_DIR)

        loader = unittest.TestLoader()
        return [next(self.test_suites_to_methods(loader.loadTestsFromName(test_name))) for test_name in test_names]

    def _run_indirect(self):
        # all available test methods
        test_runs = [TestRunDispatcher(x) for x in self.test_suites]

        test_results = []

        with concurrent.futures.ThreadPoolExecutor() as executor:
            futures = []
            for test_run in test_runs:
                f = executor.submit(
                    test_run.run
                )
                futures.append(f)

            completed_futures = concurrent.futures.as_completed(futures)

            for completed_future in completed_futures:
                test_results.append(completed_future.result())

        # print(test_results)

    def _run_direct(self):
        os.chdir(self.original_cwd)

        for test_suite in self.test_suites:
            TestRun(test_suite_to_fq_name(test_suite)).run()

    def run(self):
        if len(self.test_names) == 0:
            self.test_suites = self.get_all_test_suites()
        else:
            self.test_suites = self.test_names_to_test_suites(self.test_names)

        if self.direct:
            self._run_direct()
        else:
            self._run_indirect()



if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='System test dispatcher.')
    parser.add_argument('--tests', type=str, default=None)
    parser.add_argument('--direct', action=argparse.BooleanOptionalAction, default=False)
    args = parser.parse_args()

    test_names = []
    if not args.tests is None:
        test_names = args.tests.split(',')

    SystemTestRunner(
        args.direct,
        test_names,
        os.getcwd(),
    ).run()
