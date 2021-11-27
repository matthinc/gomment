import unittest
import os
import concurrent.futures
import argparse
import json
import subprocess
import io
import sys


SCRIPT_DIR = os.path.abspath(os.path.dirname(os.path.realpath(__file__)))


def test_suite_to_fq_name(test_suite):
    return f'{test_suite.__class__.__module__}.{test_suite.__class__.__qualname__}.{test_suite._testMethodName}'


class TestResult():
    def __init__(self, data = {}):
        self.data = data

    def __str__(self):
        return f'{self.decorate_test_status(self.get_test_status())} {self.get_test_name()}'

    def from_json(self, json_data):
        self.data = json.loads(json_data)

    def to_json(self):
        return json.dumps(self.data)

    def get_test_name(self):
        return self.data['fq_test_name']

    def get_app_output(self):
        return self.data['app_output']

    def is_successful(self):
        return self.get_test_status() == 'success'

    def get_test_status(self):
        if len(self.data['errors']) > 0:
            return 'error'
        if len(self.data['failures']) > 0:
            return 'failure'
        if self.data['tests_run'] == 1:
            return 'success'

        return 'unknown'

    def decorate_test_status(self, status):
        m = {
            'success': '\u001b[42m SUCCESS \u001b[0m',
            'failure': '\u001b[41m FAILURE \u001b[0m',
            'error':   '\u001b[45m  ERROR  \u001b[0m',
            'unknown': '\u001b[47;1m UNKNOWN \u001b[0m',
        }
        return m[status]


# class encapsulating one single test method
class TestRunDirect():
    def __init__(self, fq_test_name, json_output):
        self.fq_test_name = fq_test_name
        self.json_output = json_output

    def run(self):
        stream = sys.stderr
        stdredir = io.StringIO()

        if self.json_output:
            stream = io.StringIO()
            sys.stdout = stdredir
            sys.stderr = stdredir

        loader = unittest.TestLoader()
        suite = loader.loadTestsFromName(self.fq_test_name)

        runner = unittest.TextTestRunner(stream=stream)
        result = runner.run(suite)

        # restore original stdout + stderr
        sys.stdout = sys.__stdout__
        sys.stderr = sys.__stderr__

        if self.json_output:
            res = TestResult({
                'fq_test_name': self.fq_test_name,
                'app_output': stdredir.getvalue(),
                'test_output': stream.getvalue(),
                'failures': [x[1] for x in result.failures],
                'errors': [x[1] for x in result.errors],
                'tests_run': result.testsRun,
            })
            print(res.to_json())


# class encapsulating one single test method to be run inside a container
class TestRunIndirect():
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
          gomment-test python3 test/system_test_runner.py --direct --json --tests {self}
        '''.strip().replace('\n', ' ')

        pipe = subprocess.Popen(
            ['/bin/bash', '-c', docker_cmd],
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
        )
        stdout, stderr = pipe.communicate()
        return stdout.decode('utf-8')


class SystemTestRunner():
    def __init__(self, direct, test_names, json_output, original_cwd):
        self.direct = direct
        self.test_names = test_names
        self.test_suites = []
        self.json_output = json_output
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
        test_runs = [TestRunIndirect(x) for x in self.test_suites]

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
                res = TestResult()
                res.from_json(completed_future.result())
                test_results.append(res)

        for test_result in test_results:
            print(test_result)

    def _run_direct(self):
        os.chdir(self.original_cwd)

        for test_suite in self.test_suites:
            TestRunDirect(test_suite_to_fq_name(test_suite), self.json_output).run()

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
    parser.add_argument('--json', action=argparse.BooleanOptionalAction, default=False)
    args = parser.parse_args()

    test_names = []
    if not args.tests is None:
        test_names = args.tests.split(',')

    SystemTestRunner(
        args.direct,
        test_names,
        args.json,
        os.getcwd(),
    ).run()
