<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import Quill from 'quill';
import 'quill/dist/quill.core.css';
import 'quill/dist/quill.bubble.css';
import 'quill/dist/quill.snow.css';

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
});

const emit = defineEmits(['update:modelValue']);

const editor = ref<HTMLElement | null>(null);
let quillEditor: Quill;

onMounted(() => {
  quillEditor = new Quill(editor.value!, {
    modules: {
      toolbar: [
        [{ header: [1, 2, 3, 4, false] }],
        ['bold', 'italic', 'underline', 'strike', 'blockquote', 'code-block', 'link', 'image', 'video', 'formula'],
      ],
    },
    theme: 'snow',
    formats: ['bold', 'italic', 'underline', 'strike', 'blockquote', 'code-block', 'link', 'image', 'video', 'formula'],
    placeholder: 'Type a message',
  });
  quillEditor.root.innerHTML = props.modelValue;
  quillEditor.on('text-change', () => {
    emit('update:modelValue', quillEditor?.getText() ? quillEditor.root.innerHTML : '');
  });
});

watch(() => props.modelValue, (newValue) => {
  if (quillEditor && quillEditor.root.innerHTML !== newValue) {
    quillEditor.root.innerHTML = newValue;
  }
});
</script>
<template>
    <div ref="editor"></div>
</template>