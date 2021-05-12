import { instance, headers, plugin_host } from "./config";
import { dispatchTypicalFetch } from "./_base";

/**
 * fetchPipelines fetches the repository pipelines and dispatches an event
 * to update the store.
 */
export const fetchPipelines = (store, params) => {
    const { namespace, name,branch } = params;
  
    return dispatchTypicalFetch(store, params, "PIPELINE_LIST", () => {
      return fetch(`${instance}/v1/cicd/projects/${namespace}/${name}/pipelines?branch=${branch}`, { headers, credentials: "same-origin" });
    });
  };

/**
 * fetchPipelinesYaml fetches the repository pipelines and dispatches an event
 * to update the store.
 */
export const fetchPipelinesYaml = (store, params) => {
  const { namespace, name,branch } = params;

  return dispatchTypicalFetch(store, params, "PIPELINE_LIST_YAML", () => {
    return fetch(`${instance}/v1/cicd/projects/${namespace}/${name}/pipelines?branch=${branch}&format=yaml`, { headers, credentials: "same-origin" });
  });
};

/**
 * fetchPipeline fetches the repository pipeline and dispatches an event
 * to update the store.
 */
export const fetchPipeline =  (store, params) => {
    const { namespace, name,pipename } = params;
  
    return dispatchTypicalFetch(store, params, "PIPELINE_GET", () => {
      return fetch(`${instance}/v1/cicd/repos/pipelines/${pipename}?slug=${namespace}/${name}&format=yaml`, { headers, credentials: "same-origin" });
    });
  };


/**
 * fetchBranches fetches the repository branches and dispatches an event
 * to update the store.
 */
export const fetchBranches = (store, params) => {
    const { namespace, name } = params;
  
    return dispatchTypicalFetch(store, params, "BRANCHES_LIST", () => {
      return fetch(`${instance}/v1/cicd/repos/branches?slug=${namespace}/${name}`, { headers, credentials: "same-origin" });
    });
  };


  /**
 * createPipeline create the pipeline and dispatches an event
 * to purge the object from the store.
 */
export const createPipeline = async ({ commit }, { namespace, name, pipeline }) => {
  commit("PIPELINE_CREATE_LOADING");
  const body =pipeline.content;
  console.log(body)
  const req = await fetch(`${instance}/v1/cicd/projects/${namespace}/${name}/pipelines?branch=${pipeline.branch}`, {
    headers,
    method: "PUT",
    body,
    credentials: "same-origin"
  });

  if (req.status < 300) {
    commit("PIPELINE_CREATE_SUCCESS");
  } else {
    commit("PIPELINE_CREATE_FAILURE");
    const res = await req.json();
    throw new Error(res.message);
  }
};

  /**
 * triggerPipeline create the pipeline and dispatches an event
 * to purge the object from the store.
 */
export const triggerPipeline = async ({ commit }, { namespace, name, sign_type, sign_name}) => {
  console.log(sign_name)
  var project={
    slug: namespace+'/'+name,
    sign_type: sign_type,
    sign_name: sign_name,
  };
  const body = JSON.stringify(project);
  console.log(body)
  const req = await fetch(`${instance}/v1/cicd/builds`, {
    headers,
    method: "POST",
    body,
    credentials: "same-origin"
  });

  if (req.status < 300) {
  } else {
    const res = await req.json();
    throw new Error(res.message);
  }
};

/**
 * fetchPiplinePlugins fetches plugins and dispatches an event
 * to update the store.
 */
export const fetchPipelinePlugins = (store, params) => {
 
  return dispatchTypicalFetch(store, params, "PLUGINS_LIST", () => {
    return fetch(`${instance}/v1/cicd/plugins`, { headers, credentials: "same-origin" });
  });
};

/**
 * fetchPiplinePlugins fetches plugins and dispatches an event
 * to update the store.
 */
export const fetchPluginsVersions = (store, params) => {
  const {name } = params;
  return dispatchTypicalFetch(store, params, "PLUGINS_VERSIONS", () => {
    return fetch(`${instance}/v1/cicd/plugins/${name}/versions`, { headers, credentials: "same-origin" });
  });
};
