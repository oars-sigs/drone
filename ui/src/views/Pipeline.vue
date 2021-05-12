<template>
  <CardGroup class="settings-view">
    <Card contentPadding="0 15px 15px">
      <div class="control-group">
        <label class="control-label">Branch</label>
        <div class="controls">
          <BaseSelect :value="''" @input="handleGet" :options="branches" />
        </div>
        <Button v-if="curentBranch" outline @click.native="trigger" :loading="triggering">Trigger</Button>
      </div>

      <div v-if="ready" class="control-group">
        <Editor
          class="editor"
          v-model="pipeline.content"
          @init="editorInit"
          lang="yaml"
          theme="chrome"
          height="300"
        ></Editor>
      </div>
      <div v-if="ready" class="control-actions">
        <Button theme="primary" v-if="isAdmin" size="l" @click.native="save" :loading="saving">Save</Button>
        <div class="error-message" v-if="error">{{ error.message }}</div>
        
      </div>
    </Card>
  </CardGroup>
</template>

<script>
import Editor from "vue2-ace-editor";
import BaseSelect from "@/components/forms/BaseSelect.vue";
import Card from "@/components/Card.vue";
import CardGroup from "@/components/CardGroup.vue";
import Button from "@/components/buttons/Button.vue";
import { instance, headers } from "@/actions/config";

export default {
  name: "pipelines",
  data() {
    return {
      name: "",
      ready: false,
      pipeline: {
        content: "",
        branch:""
      },
      branches:[],
      curentBranch:"",
      error: null,
      saving: false,
      trgError: null,
      namespace: '',
      repoName:'',
      triggering: false
    };
  },
  components: {
    BaseSelect,
    Card,
    CardGroup,
    Button,
    Editor,
  },
  created() {
      this.namespace=this.$route.params.namespace;
      this.repoName=this.$route.params.name;
      this.listBranches()
      this.handleGet('')
      
  },
  computed: {
    slug() {
      return this.$route.params.namespace + "/" + this.$route.params.name;
    },
    repo() {
      let repo = this.$store.state.repos[this.slug];
      return repo && { ...repo };
    },
    isRoot() {
      return this.$store.state.user.data.admin;
    },
    isAdmin() {
      const isAdmin = this.repo && this.repo.permissions && this.repo.permissions.admin;
      console.log(isAdmin)
      return this.isRoot || isAdmin;
    },
  },
  methods: {
    editorInit: function() {
      require("brace/ext/language_tools"); //language extension prerequsite...
      require("brace/theme/chrome");
      require("brace/mode/yaml");
    },
    handleGet(value) {
      this.curentBranch=value;
      var _that =this
      fetch(`/extend/${this.namespace}/${this.repoName}/pipelines?branch=${value}`, {headers}).then(function(response) {
        return response.json()
      }).then(function(json) {
        _that.pipeline.content=json.data
        _that.ready=true
      }).catch(function(ex) {
        console.log('parsing failed', ex) 
      })
    },
    listBranches(){
      var _that =this
      fetch(`/extend/${this.namespace}/${this.repoName}/branches`, 
      {
        headers,
      }).then(function(response) {
        return response.json()
      }).then(function(json) {
        _that.branches=[["","默认"]]
        json.forEach(element => {
          _that.branches.push([element.Name,element.Name])
        });
        
      }).catch(function(ex) {
        console.log('parsing failed', ex) 
      })
    },
    save() {
      var _that =this
      this.saving = true;
       fetch(`/extend/${this.namespace}/${this.repoName}/pipelines?branch=${this.curentBranch}`, 
      {
        method: "PUT",
        headers,
        body: this.pipeline.content,
      }).then(function(response) {
          _that.$store.dispatch("showNotification", { message: "Successfully saved" });
          _that.error = null;
          _that.saving = false;
      }).catch(function(ex) {
        _that.error = ex;
      })
    },
    trigger(){
      var _that =this
      this.triggering = true;
       fetch(`/extend/${this.namespace}/${this.repoName}/builds?branch=${this.curentBranch}`, 
      {
        method: "POST",
        headers,
      }).then(function(response) {
          _that.$store.dispatch("showNotification", { message: "Successfully trigger" });
          _that.triggering = false;
      }).catch(function(ex) {
        _that.$store.dispatch("showNotification", { message: ex });
      })
    }
  }
};
</script>

<style scoped lang="scss">
@import "../assets/styles/mixins";
@import "../assets/styles/mixins";

.control-group {
  .controls {
    & + .help {
      flex-shrink: 0;
    }
    .base-checkbox + .base-checkbox {
      margin-left: 48px;
      @include tablet {
        margin-left: 0px;
        margin-top: 10px;
        display: block;
      }
    }
  }
  @include tablet {
    flex-direction: row;
    flex-wrap: wrap;
  }
}
.control-label {
  @include tablet {
    line-height: 18px;
    order: -2;
    flex-grow: 1;
  }
}
.help {
  @include tablet {
    order: -1;
  }
}
.help-p + .help-p {
  margin-top: 10px;
}

.disable {
  padding: 0 15px;

  @include mobile {
    text-align: center;
  }
}

.disable span {
  margin: 15px 0 0 15px;
  color: $color-text-secondary;

  @include mobile {
    display: block;
    margin: 10px 0 0 0;
  }
}
.editor {
  font-size: 15px;
}

// .trigger{
//   margin-left: 50px;
// }
</style>
