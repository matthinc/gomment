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
        fqname = self.data['fq_test_name']
        return f'TestResult({fqname})'

    def get_decorated_string(self):
        return f'{self.decorate_test_status(self.get_test_status())} {self.get_test_name()}'

    def from_json(self, json_data):
        self.data = json.loads(json_data)

    def to_json(self):
        return json.dumps(self.data)

    def get_test_name(self):
        return self.data['fq_test_name']

    def get_app_output(self):
        return self.data['app_output']

    def get_test_output(self):
        return self.data['test_output']

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
    def __init__(self, test_suite, passthrough):
        self.test_suite = test_suite
        self.passthrough = passthrough

    def __str__(self):
        return test_suite_to_fq_name(self.test_suite)

    def run(self):
        state_dir = f'{SCRIPT_DIR}/state/{self}'

        if not os.path.exists(state_dir):
            os.makedirs(state_dir)

        passthrough_str = ''
        if self.passthrough:
            passthrough_str = '-v$(pwd)/gomment:/app/gomment:ro'

        docker_cmd = f'''
        docker run
          --rm
          -u$(id -u):$(id -g)
          -v{SCRIPT_DIR}:/app/test
          -v{state_dir}:/app/test-state
          {passthrough_str}
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


class TestDiscovery():
    """Converts a set of test names into a set of executable test
    suites."""

    def __init__(self, test_names):
        self.test_names = test_names

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

    def get_test_suites_by_name(self, test_name):
        os.chdir(SCRIPT_DIR)
        loader = unittest.TestLoader()

        tests = loader.loadTestsFromName(test_name)

        return self.test_suites_to_methods(tests)

    def get_test_suites(self):
        if len(self.test_names) == 0:
            return self.get_all_test_suites()

        test_suites = []
        for test_name in self.test_names:
            test_suites.extend(self.get_test_suites_by_name(test_name))

        return test_suites


class SystemTestRunner():
    def __init__(self, direct, test_names, json_output, passthrough, original_cwd):
        self.direct = direct
        self.test_names = test_names
        self.test_suites = []
        self.json_output = json_output
        self.passthrough = passthrough
        self.original_cwd = original_cwd

    def _run_indirect(self):
        # all available test methods
        test_runs = [TestRunIndirect(x, self.passthrough) for x in self.test_suites]

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
                try:
                    res.from_json(completed_future.result())
                    test_results.append(res)
                except:
                    print(completed_future.result())

        test_results.sort(key=lambda x: str(x))

        # print all unsuccessful traces
        for test_result in test_results:
            if not test_result.is_successful():
                print(test_result.get_decorated_string())
                print(test_result.get_app_output())
                print(test_result.get_test_output())

        # print summary
        for test_result in test_results:
            print(test_result.get_decorated_string())

    def _run_direct(self):
        for test_suite in self.test_suites:
            TestRunDirect(test_suite_to_fq_name(test_suite), self.json_output).run()

    def run(self):
        self.test_suites = TestDiscovery(self.test_names).get_test_suites()

        os.chdir(self.original_cwd)

        if self.direct:
            self._run_direct()
        else:
            self._run_indirect()



if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='System test dispatcher.')
    parser.add_argument('--tests', type=str, default=None)
    parser.add_argument('--direct', action=argparse.BooleanOptionalAction, default=False)
    parser.add_argument('--json', action=argparse.BooleanOptionalAction, default=False)
    parser.add_argument('--passthrough', action=argparse.BooleanOptionalAction, default=False, help='Enable passthrough for gomment binary. Enables quick testing without recompiling the docker image. Requires compatible glibc.')
    args = parser.parse_args()

    if args.passthrough and args.direct:
        print("arguments --passthrough and --direct are incompatible")
        sys.exit(2)

    test_names = []
    if not args.tests is None:
        test_names = args.tests.split(',')

    SystemTestRunner(
        args.direct,
        test_names,
        args.json,
        args.passthrough,
        os.getcwd(),
    ).run()
