import json
import time
import os
import argparse
import tarfile
import requests


fateflow_uri = os.getenv('FateFlowServer')
server_url = "http://{}/{}".format(fateflow_uri, "v1")
print(server_url)


def submit_job(dsl, config):
    post_data = {'job_dsl': dsl,
                 'job_runtime_conf': config}
    response = requests.post(
        "/".join([server_url, "job", "submit"]), json=post_data)

    return prettify(response)


def query_job(query_conditions):
    response = requests.post(
        "/".join([server_url, "job", "query"]), json=query_conditions)
    return prettify(response)


def stop_job(self, job_id):
    post_data = {
        'job_id': job_id
    }
    response = requests.post(
        "/".join([server_url, "job", "stop"]), json=post_data)
    return prettify(response)


def prettify(response, verbose=False):
    if verbose:
        if isinstance(response, requests.Response):
            if response.status_code == 200:
                print("Success!")
            print(json.dumps(response.json(), indent=4, ensure_ascii=False))
        else:
            print(response)

    return response


def run_job(job_dsl, config_data):
    print(job_dsl)
    print(config_data)
    response = submit_job(job_dsl, config_data)
    prettify(response, True)
    stdout = json.loads(response.content)
    print(stdout)
    jobid = stdout["jobId"]
    query_condition = {
        "job_id": jobid
    }
    job_status = query_job(query_condition)
    prettify(job_status, True)

    for i in range(500):
        time.sleep(1)
        job_detail = query_job(query_condition).json()
        final_status = job_detail["data"][0]["f_status"]
        print(final_status)

        if final_status == "failed":
            print("Failed")
            break
        if final_status == "success":
            print("Success")
            break

        # r
        # esponse = manager.fetch_job_log_new(jobid)


if __name__ == "__main__":
    arg_parser = argparse.ArgumentParser()
    arg_parser.add_argument("--dsl", type=str, help="please input dsl")
    arg_parser.add_argument("--config", type=str, help="please input config")
    args = arg_parser.parse_args()
    print(type(args.dsl))
    print(args.config)
    dsl = json.loads(args.dsl)
    config = json.loads(args.config)

    run_job(dsl, config)
    exit(0)
