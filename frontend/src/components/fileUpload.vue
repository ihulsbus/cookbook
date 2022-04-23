<template>
  <FileUpload
    mode="advanced"
    name="multiplefiles"
    :maxFileSize="10000000"
    accept="image/*"
    :multiple="true"
    @upload="onUploadSuccess"
    @error="onUploadError"
    @before-send="beforeSend"
  />
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import FileUpload from 'primevue/fileupload';
import idsrvAuth from '@/lib/auth/auth';

interface RequestObject {
    xhr: XMLHttpRequest;
    formData: FormData;
}

@Options({
  components: {
    FileUpload,
  },
  methods: {
    beforeSend(request:RequestObject) {
      request.xhr.open('POST', `http://localhost:8080/v1/recipe/${this.$route.params.id}/upload`);
      request.xhr.setRequestHeader('Authorization', `Bearer ${idsrvAuth.accessToken}`);
      return request;
    },
    onUploadSuccess() {
      this.$toast.add({
        severity: 'success', summary: 'Success', detail: 'File uploaded', life: 3000,
      });
    },
    onUploadError() {
      this.$toast.add({
        severity: 'error', summary: 'Error', detail: 'File upload failed', life: 3000,
      });
    },
  },
})
export default class fileUpload extends Vue {}
</script>

<style lang="scss">
  .surface-ground {
    height: 100%
  }
</style>
