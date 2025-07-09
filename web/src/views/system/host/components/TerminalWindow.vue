<template>
  <a-modal
    v-model:open="isVisible"
    title="WebSSH Terminal"
    :footer="null"
    @cancel="handleCancel"
    width="80vw"
    wrap-class-name="full-height-modal"
  >
    <div ref="terminalRef" class="terminal-container"></div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue';
import type { PropType } from 'vue';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import 'xterm/css/xterm.css';
import type { Host } from '@/types/api/host';
import { getAuthToken } from '@/utils/auth'; // Assuming you have such a utility
import { useUserStore } from '@/store/modules/user';

const props = defineProps({
  visible: { type: Boolean, required: true },
  host: { type: Object as PropType<Host | null>, required: true },
});

const emit = defineEmits(['update:visible']);

const userStore = useUserStore();

const terminalRef = ref<HTMLElement | null>(null);
const isVisible = ref(false);
let term: Terminal;
let ws: WebSocket;
let fitAddon: FitAddon;

const termOptions = {
  cursorBlink: true,
  fontSize: 14,
  fontFamily: 'Consolas, "Courier New", monospace',
  theme: {
    background: '#1e1e1e',
    foreground: '#d4d4d4',
  },
};

const setupTerminal = () => {
  if (!terminalRef.value || !props.host) return;

  term = new Terminal(termOptions);
  fitAddon = new FitAddon();
  term.loadAddon(fitAddon);

  term.open(terminalRef.value);
  fitAddon.fit();

  const token = userStore.accessToken; // Use the store directly
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const wsUrl = `${protocol}//${window.location.host}/api/v1/cmdb/ws/ssh?host_id=${props.host.id}&auth=${token}`;

  ws = new WebSocket(wsUrl);

  ws.onopen = () => {
    term.focus();
    fitAddon.fit();
  };

  ws.onmessage = (event) => {
    term.write(event.data);
  };

  ws.onclose = () => {
    term.write('\r\nConnection closed.\r\n');
  };

  ws.onerror = (error) => {
    console.error('WebSocket Error:', error);
    term.write('\r\nAn error occurred.\r\n');
  };

  term.onData((data) => {
    ws.send(data);
  });

  window.addEventListener('resize', () => fitAddon.fit());
};

const disposeTerminal = () => {
  if (ws) {
    ws.close();
  }
  if (term) {
    term.dispose();
  }
  window.removeEventListener('resize', () => fitAddon.fit());
};

watch(() => props.visible, (newValue) => {
  isVisible.value = newValue;
  if (newValue) {
    nextTick(() => {
      setupTerminal();
    });
  } else {
    disposeTerminal();
  }
});

onUnmounted(() => {
  disposeTerminal();
});

const handleCancel = () => {
  emit('update:visible', false);
};
</script>

<style lang="scss">
.terminal-container {
  width: 100%;
  height: 60vh;
}
.full-height-modal .ant-modal {
  max-width: 100%;
  top: 0;
  padding-bottom: 0;
  margin: 0;
}
.full-height-modal .ant-modal-content {
  display: flex;
  flex-direction: column;
  height: calc(100vh);
}
.full-height-modal .ant-modal-body {
  flex: 1;
}
</style> 