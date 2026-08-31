<template>
  <div class="page-container">
    <div class="card">
      <div class="table-head">
        <h3 style="margin: 0">节点分组</h3>
        <t-space size="small">
          <t-button theme="primary" @click="openGroupCreate">新增分组</t-button>
          <t-button theme="default" variant="outline" @click="loadGroups">刷新</t-button>
        </t-space>
      </div>
      <t-table :data="groups" :columns="groupColumns" row-key="id" :loading="groupLoading">
        <template #op="{ row }">
          <t-space size="small">
            <t-button theme="primary" variant="text" @click="openGroupEdit(row)">编辑</t-button>
            <t-button theme="danger" variant="text" @click="onDropGroup(row)">删除</t-button>
          </t-space>
        </template>
      </t-table>
    </div>

    <div class="card" style="margin-top: 16px">
      <div class="table-head">
        <h3 style="margin: 0">节点管理</h3>
        <t-space size="small">
          <t-button theme="primary" @click="openNodeCreate">新增节点</t-button>
          <t-button theme="default" variant="outline" @click="loadNodes">刷新</t-button>
        </t-space>
      </div>
      <t-tabs v-model="nodeType" @change="loadNodes">
        <t-tab-panel v-for="t in nodeTypes" :key="t.value" :value="t.value" :label="t.label" />
      </t-tabs>
      <t-table :data="nodes" :columns="nodeColumns" row-key="id" :loading="nodeLoading">
        <template #op="{ row }">
          <t-space size="small">
            <t-button theme="primary" variant="text" @click="openNodeEdit(row)">编辑</t-button>
            <t-button theme="success" variant="text" @click="openInstall(row)">安装</t-button>
            <t-button theme="danger" variant="text" @click="onDropNode(row)">删除</t-button>
          </t-space>
        </template>
      </t-table>
    </div>

    <t-dialog v-model:visible="groupDialogVisible" :header="groupForm.id ? '编辑分组' : '新增分组'" :confirm-btn="{ onClick: submitGroup }" :on-close="resetGroupForm">
      <t-form :data="groupForm" label-width="96px">
        <t-form-item label="名称" required>
          <t-input v-model="groupForm.name" placeholder="分组名称" />
        </t-form-item>
      </t-form>
    </t-dialog>

    <t-dialog v-model:visible="nodeDialogVisible" :header="nodeForm.id ? '编辑节点' : '新增节点'" width="560px" :confirm-btn="{ onClick: submitNode }" :on-close="resetNodeForm">
      <t-form :data="nodeForm" label-width="110px">
        <t-form-item label="名称" required>
          <t-input v-model="nodeForm.name" placeholder="节点名称" />
        </t-form-item>
        <t-form-item label="分组" required>
          <t-select v-model="nodeForm.group_id" placeholder="选择分组" clearable filterable allow-input>
            <t-option v-for="g in groups" :key="g.id" :value="String(g.id)" :label="g.name" />
          </t-select>
        </t-form-item>
        <t-form-item label="地址" required>
          <t-input v-model="nodeForm.host" placeholder="主机地址" />
        </t-form-item>
        <t-form-item label="端口" required>
          <t-input v-model="nodeForm.port" placeholder="连接端口" />
        </t-form-item>
        <t-form-item label="服务端口">
          <t-input-number v-model="nodeForm.server_port" :min="0" :step="1" />
        </t-form-item>
        <t-form-item label="倍率">
          <t-input v-model="nodeForm.rate" placeholder="1" />
        </t-form-item>
        <template v-if="nodeType === 'trojan'">
          <t-form-item label="允许不安全">
            <t-switch v-model="nodeForm.allow_insecure" />
          </t-form-item>
          <t-form-item label="SNI">
            <t-input v-model="nodeForm.server_name" placeholder="server_name" />
          </t-form-item>
        </template>
        <template v-else-if="nodeType === 'vmess'">
          <t-form-item label="TLS">
            <t-select v-model="nodeForm.tls" :options="[{ label: '关闭', value: 0 }, { label: '开启', value: 1 }]" />
          </t-form-item>
          <t-form-item label="传输">
            <t-select v-model="nodeForm.network" :options="['tcp', 'ws', 'grpc', 'http', 'kcp'].map((v) => ({ label: v, value: v }))" />
          </t-form-item>
        </template>
        <template v-else-if="nodeType === 'shadowsocks'">
          <t-form-item label="加密">
            <t-select v-model="nodeForm.cipher" :options="['aes-256-gcm', 'aes-128-gcm', 'chacha20-ietf-poly1305'].map((v) => ({ label: v, value: v }))" />
          </t-form-item>
        </template>
        <template v-else-if="nodeType === 'hysteria'">
          <t-form-item label="上行(Mbps)">
            <t-input-number v-model="nodeForm.up_mbps" :min="0" :step="1" />
          </t-form-item>
          <t-form-item label="下行(Mbps)">
            <t-input-number v-model="nodeForm.down_mbps" :min="0" :step="1" />
          </t-form-item>
          <t-form-item label="SNI">
            <t-input v-model="nodeForm.server_name" placeholder="server_name" />
          </t-form-item>
          <t-form-item label="允许不安全">
            <t-switch v-model="nodeForm.insecure" />
          </t-form-item>
        </template>
        <t-form-item label="排序">
          <t-input-number v-model="nodeForm.sort" :min="0" :step="1" />
        </t-form-item>
        <t-form-item label="展示">
          <t-switch v-model="nodeForm.show" />
        </t-form-item>
      </t-form>
    </t-dialog>

    <t-dialog v-model:visible="installDialogVisible" header="远程安装节点" width="560px" :confirm-btn="{ loading: installing, onClick: submitInstall }">
      <t-form :data="installForm" label-width="110px">
        <t-form-item label="SSH 主机" required>
          <t-input v-model="installForm.ssh_host" placeholder="节点 SSH 地址" />
        </t-form-item>
        <t-form-item label="SSH 端口">
          <t-input-number v-model="installForm.ssh_port" :min="1" :max="65535" />
        </t-form-item>
        <t-form-item label="用户名" required>
          <t-input v-model="installForm.ssh_user" placeholder="root" />
        </t-form-item>
        <t-form-item label="密码">
          <t-input v-model="installForm.ssh_password" type="password" placeholder="与私钥二选一" />
        </t-form-item>
        <t-form-item label="私钥">
          <t-textarea v-model="installForm.ssh_private_key" :autosize="{ minRows: 4, maxRows: 8 }" placeholder="SSH 私钥内容（可选）" />
        </t-form-item>
      </t-form>
    </t-dialog>

    <t-dialog v-model:visible="installLogVisible" header="安装日志" width="640px" :footer="false">
      <t-textarea :model-value="installLogText" readonly :autosize="{ minRows: 14, maxRows: 24 }" />
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, reactive, ref, computed, h } from 'vue';
import { MessagePlugin, Tag } from 'tdesign-vue-next';
import {
  listGroups,
  saveGroup,
  dropGroup,
  listNodes,
  saveNode,
  dropNode,
  installNode,
  nodeInstallLogUrl,
  ServerGroup,
  AdminServerGroupSaveReq,
  ServerNode,
  AdminServerNodeSaveReq,
} from '@/api/admin';

const nodeTypes = [
  { label: 'Trojan', value: 'trojan' },
  { label: 'Vmess', value: 'vmess' },
  { label: 'Shadowsocks', value: 'shadowsocks' },
  { label: 'Hysteria', value: 'hysteria' },
];

const groups = ref<ServerGroup[]>([]);
const groupLoading = ref(false);
const groupDialogVisible = ref(false);
const groupForm = reactive<AdminServerGroupSaveReq>({ id: undefined, name: '' });

const nodes = ref<ServerNode[]>([]);
const nodeLoading = ref(false);
const nodeType = ref<string>('trojan');
const nodeDialogVisible = ref(false);
const nodeForm = reactive<AdminServerNodeSaveReq>({
  id: undefined,
  type: '',
  name: '',
  group_id: '',
  host: '',
  port: '',
  server_port: 0,
  rate: '1',
  show: true,
  sort: 0,
});

const groupColumns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'name', title: '名称' },
  { colKey: 'created_at', title: '创建时间', cell: (_: unknown, row: ServerGroup) => formatTime(row.created_at) },
  { colKey: 'op', title: '操作', width: 140 },
];

const nodeColumns = computed(() => [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'name', title: '名称' },
  { colKey: 'group_id', title: '分组', width: 90 },
  { colKey: 'host', title: '地址' },
  { colKey: 'port', title: '端口', width: 90 },
  { colKey: 'server_port', title: '服务端口', width: 100 },
  { colKey: 'rate', title: '倍率', width: 80 },
  { colKey: 'show', title: '展示', cell: (_: unknown, row: ServerNode) => (row.show ? '是' : '否') },
  { colKey: 'online', title: '状态', width: 110, cell: (_: unknown, row: ServerNode) => renderStatus(row) },
  { colKey: 'op', title: '操作', width: 200 },
]);

async function loadGroups() {
  groupLoading.value = true;
  try {
    groups.value = await listGroups();
  } finally {
    groupLoading.value = false;
  }
}

function resetGroupForm() {
  Object.assign(groupForm, { id: undefined, name: '' });
}

function openGroupCreate() {
  resetGroupForm();
  groupDialogVisible.value = true;
}

function openGroupEdit(row: ServerGroup) {
  resetGroupForm();
  groupForm.id = row.id;
  groupForm.name = row.name;
  groupDialogVisible.value = true;
}

async function submitGroup() {
  if (!groupForm.name) {
    MessagePlugin.error('请输入分组名称');
    return;
  }
  await saveGroup({ id: groupForm.id, name: groupForm.name });
  MessagePlugin.success('已保存');
  groupDialogVisible.value = false;
  loadGroups();
}

async function onDropGroup(row: ServerGroup) {
  await dropGroup(row.id);
  MessagePlugin.success('已删除');
  loadGroups();
}

async function loadNodes() {
  nodeLoading.value = true;
  try {
    nodes.value = await listNodes(nodeType.value);
  } finally {
    nodeLoading.value = false;
  }
}

function resetNodeForm() {
  Object.assign(nodeForm, {
    id: undefined,
    name: '',
    group_id: '',
    host: '',
    port: '',
    server_port: 0,
    rate: '1',
    show: true,
    sort: 0,
    allow_insecure: false,
    server_name: '',
    tls: 0,
    network: 'tcp',
    cipher: 'aes-256-gcm',
    up_mbps: 0,
    down_mbps: 0,
    insecure: false,
  });
}

function openNodeCreate() {
  resetNodeForm();
  nodeDialogVisible.value = true;
}

function openNodeEdit(row: ServerNode) {
  resetNodeForm();
  Object.assign(nodeForm, {
    id: row.id,
    name: row.name,
    group_id: row.group_id,
    host: row.host,
    port: row.port,
    server_port: row.server_port,
    rate: row.rate,
    show: row.show,
    sort: row.sort ?? 0,
    allow_insecure: row.allow_insecure ?? false,
    server_name: row.server_name ?? '',
    tls: row.tls ?? 0,
    network: row.network ?? 'tcp',
    cipher: row.cipher ?? 'aes-256-gcm',
    up_mbps: row.up_mbps ?? 0,
    down_mbps: row.down_mbps ?? 0,
    insecure: row.insecure ?? false,
  });
  nodeDialogVisible.value = true;
}

async function submitNode() {
  if (!nodeForm.name || !nodeForm.host || !nodeForm.port) {
    MessagePlugin.error('请填写完整节点信息');
    return;
  }
  const data: AdminServerNodeSaveReq = {
    id: nodeForm.id,
    type: nodeType.value,
    name: nodeForm.name,
    group_id: nodeForm.group_id || '1',
    host: nodeForm.host,
    port: nodeForm.port,
    server_port: nodeForm.server_port ?? 0,
    rate: nodeForm.rate || '1',
    show: nodeForm.show,
    sort: nodeForm.sort,
    allow_insecure: nodeForm.allow_insecure,
    server_name: nodeForm.server_name || undefined,
    tls: nodeForm.tls,
    network: nodeForm.network || undefined,
    cipher: nodeForm.cipher || undefined,
    up_mbps: nodeForm.up_mbps,
    down_mbps: nodeForm.down_mbps,
    insecure: nodeForm.insecure,
  };
  await saveNode(data);
  MessagePlugin.success('已保存');
  nodeDialogVisible.value = false;
  loadNodes();
}

async function onDropNode(row: ServerNode) {
  await dropNode(nodeType.value, row.id);
  MessagePlugin.success('已删除');
  loadNodes();
}

function formatTime(v?: number): string {
  if (!v) return '-';
  return new Date(v * 1000).toLocaleString();
}

function renderStatus(row: ServerNode) {
  if (row.online) {
    return h(Tag, { theme: 'success', variant: 'light' }, () => '在线');
  }
  return h(Tag, { theme: 'default', variant: 'outline' }, () => '离线');
}

const installDialogVisible = ref(false);
const installLogVisible = ref(false);
const installForm = reactive({ ssh_host: '', ssh_port: 22, ssh_user: 'root', ssh_password: '', ssh_private_key: '' });
const installLogText = ref('');
const installing = ref(false);
let currentNode: ServerNode | null = null;
let es: EventSource | null = null;

function openInstall(row: ServerNode) {
  currentNode = row;
  installForm.ssh_host = row.host;
  installForm.ssh_port = 22;
  installForm.ssh_user = 'root';
  installForm.ssh_password = '';
  installForm.ssh_private_key = '';
  installLogText.value = '';
  installDialogVisible.value = true;
}

async function submitInstall() {
  if (!currentNode) return;
  if (!installForm.ssh_host || !installForm.ssh_user) {
    MessagePlugin.error('请填写 SSH 信息');
    return;
  }
  installing.value = true;
  try {
    const { task_id } = await installNode({
      id: currentNode.id,
      type: nodeType.value,
      ssh_host: installForm.ssh_host,
      ssh_port: installForm.ssh_port,
      ssh_user: installForm.ssh_user,
      ssh_password: installForm.ssh_password || undefined,
      ssh_private_key: installForm.ssh_private_key || undefined,
    });
    installDialogVisible.value = false;
    openLog(task_id);
  } catch {
    MessagePlugin.error('发起安装失败');
  } finally {
    installing.value = false;
  }
}

function openLog(taskID: string) {
  installLogText.value = '';
  installLogVisible.value = true;
  es?.close();
  es = new EventSource(nodeInstallLogUrl(taskID));
  es.onmessage = (ev) => {
    installLogText.value += ev.data + '\n';
    if (ev.data === '##INSTALL_DONE##' || ev.data.startsWith('##INSTALL_FAILED')) {
      es?.close();
      MessagePlugin.success(ev.data === '##INSTALL_DONE##' ? '安装完成' : '安装失败');
      loadNodes();
    }
  };
  es.onerror = () => es?.close();
}

onUnmounted(() => es?.close());

onMounted(() => {
  loadGroups();
  loadNodes();
});
</script>

<style scoped>
.table-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
</style>
