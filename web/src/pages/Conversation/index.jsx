// web/src/pages/Conversation/index.jsx
// 对话记录管理页面

import React, { useEffect, useState } from 'react';
import { Button, Table, Form, Modal, DatePicker, Input, Space, Tag, Popconfirm } from '@douyinfe/semi-ui';
import { API, showError, showSuccess, showInfo } from '../../helpers';
import { renderTimestamp } from '../../helpers/render';

const ConversationManagement = () => {
  const [conversations, setConversations] = useState([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [selectedKeys, setSelectedKeys] = useState([]);
  const [detailVisible, setDetailVisible] = useState(false);
  const [currentDetail, setCurrentDetail] = useState(null);
  const [settingVisible, setSettingVisible] = useState(false);
  const [conversationLogEnabled, setConversationLogEnabled] = useState(false);

  // 筛选条件
  const [filters, setFilters] = useState({
    username: '',
    model_name: '',
    start_time: 0,
    end_time: 0,
  });

  useEffect(() => {
    loadConversations();
    loadSetting();
  }, [page, pageSize]);

  // 加载对话记录列表
  const loadConversations = async () => {
    setLoading(true);
    try {
      const params = {
        page,
        page_size: pageSize,
        ...filters,
      };
      const res = await API.get('/api/conversation/', { params });
      if (res.data.success) {
        setConversations(res.data.data.data || []);
        setTotal(res.data.data.total || 0);
      } else {
        showError('加载失败：' + res.data.message);
      }
    } catch (error) {
      showError('加载失败：' + error.message);
    } finally {
      setLoading(false);
    }
  };

  // 加载功能开关设置
  const loadSetting = async () => {
    try {
      const res = await API.get('/api/conversation/setting');
      if (res.data.success) {
        setConversationLogEnabled(res.data.data.enabled);
      }
    } catch (error) {
      console.error('加载设置失败', error);
    }
  };

  // 更新功能开关
  const updateSetting = async (enabled) => {
    try {
      const res = await API.put('/api/conversation/setting', { enabled });
      if (res.data.success) {
        showSuccess('设置已更新');
        setConversationLogEnabled(enabled);
        setSettingVisible(false);
      } else {
        showError('更新失败：' + res.data.message);
      }
    } catch (error) {
      showError('更新失败：' + error.message);
    }
  };

  // 查看详情
  const viewDetail = async (id) => {
    try {
      const res = await API.get(`/api/conversation/${id}`);
      if (res.data.success) {
        setCurrentDetail(res.data.data);
        setDetailVisible(true);
      } else {
        showError('加载失败：' + res.data.message);
      }
    } catch (error) {
      showError('加载失败：' + error.message);
    }
  };

  // 批量删除
  const batchDelete = async () => {
    if (selectedKeys.length === 0) {
      showInfo('请选择要删除的记录');
      return;
    }

    try {
      const res = await API.delete('/api/conversation/', {
        data: { ids: selectedKeys },
      });
      if (res.data.success) {
        showSuccess(`成功删除 ${res.data.data.deleted} 条记录`);
        setSelectedKeys([]);
        loadConversations();
      } else {
        showError('删除失败：' + res.data.message);
      }
    } catch (error) {
      showError('删除失败：' + error.message);
    }
  };

  // 按条件删除
  const deleteByCondition = async () => {
    if (!filters.username && !filters.model_name && !filters.start_time && !filters.end_time) {
      showInfo('请至少选择一个筛选条件');
      return;
    }

    try {
      const res = await API.post('/api/conversation/delete_by_condition', filters);
      if (res.data.success) {
        showSuccess(`成功删除 ${res.data.data.deleted} 条记录`);
        loadConversations();
      } else {
        showError('删除失败：' + res.data.message);
      }
    } catch (error) {
      showError('删除失败：' + error.message);
    }
  };

  // 搜索
  const handleSearch = () => {
    setPage(1);
    loadConversations();
  };

  // 重置筛选
  const handleReset = () => {
    setFilters({
      username: '',
      model_name: '',
      start_time: 0,
      end_time: 0,
    });
    setPage(1);
    setTimeout(() => loadConversations(), 100);
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      width: 120,
    },
    {
      title: '模型',
      dataIndex: 'model_name',
      width: 150,
      render: (text) => <Tag>{text}</Tag>,
    },
    {
      title: 'Token名称',
      dataIndex: 'token_name',
      width: 150,
    },
    {
      title: '输入Token',
      dataIndex: 'prompt_tokens',
      width: 100,
    },
    {
      title: '输出Token',
      dataIndex: 'completion_tokens',
      width: 100,
    },
    {
      title: '总Token',
      dataIndex: 'total_tokens',
      width: 100,
    },
    {
      title: '是否流式',
      dataIndex: 'is_stream',
      width: 100,
      render: (isStream) => (isStream ? '是' : '否'),
    },
    {
      title: '响应时间(ms)',
      dataIndex: 'use_time',
      width: 120,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      width: 180,
      render: (text) => renderTimestamp(text),
    },
    {
      title: '操作',
      fixed: 'right',
      width: 150,
      render: (_, record) => (
        <Space>
          <Button size="small" onClick={() => viewDetail(record.id)}>
            查看详情
          </Button>
          <Popconfirm
            title="确定删除吗？"
            onConfirm={() => {
              API.delete('/api/conversation/', { data: { ids: [record.id] } }).then((res) => {
                if (res.data.success) {
                  showSuccess('删除成功');
                  loadConversations();
                } else {
                  showError('删除失败：' + res.data.message);
                }
              });
            }}
          >
            <Button size="small" type="danger">
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: 24 }}>
      <h2>对话记录管理</h2>

      {/* 功能开关状态 */}
      <div style={{ marginBottom: 16, padding: 12, background: '#f6f7f9', borderRadius: 4 }}>
        <Space>
          <span>对话记录功能：</span>
          <Tag color={conversationLogEnabled ? 'green' : 'red'}>
            {conversationLogEnabled ? '已启用' : '已禁用'}
          </Tag>
          <Button size="small" onClick={() => setSettingVisible(true)}>
            修改设置
          </Button>
        </Space>
      </div>

      {/* 筛选表单 */}
      <Form layout="horizontal" style={{ marginBottom: 16 }}>
        <Space>
          <Input
            placeholder="用户名"
            value={filters.username}
            onChange={(value) => setFilters({ ...filters, username: value })}
            style={{ width: 150 }}
          />
          <Input
            placeholder="模型名称"
            value={filters.model_name}
            onChange={(value) => setFilters({ ...filters, model_name: value })}
            style={{ width: 150 }}
          />
          <DatePicker
            type="dateTimeRange"
            placeholder={['开始时间', '结束时间']}
            onChange={(value) => {
              if (value && value.length === 2) {
                setFilters({
                  ...filters,
                  start_time: Math.floor(value[0].getTime() / 1000),
                  end_time: Math.floor(value[1].getTime() / 1000),
                });
              } else {
                setFilters({ ...filters, start_time: 0, end_time: 0 });
              }
            }}
            style={{ width: 300 }}
          />
          <Button type="primary" onClick={handleSearch}>
            搜索
          </Button>
          <Button onClick={handleReset}>重置</Button>
        </Space>
      </Form>

      {/* 批量操作 */}
      <div style={{ marginBottom: 16 }}>
        <Space>
          <Popconfirm
            title={`确定删除选中的 ${selectedKeys.length} 条记录吗？`}
            disabled={selectedKeys.length === 0}
            onConfirm={batchDelete}
          >
            <Button type="danger" disabled={selectedKeys.length === 0}>
              批量删除 ({selectedKeys.length})
            </Button>
          </Popconfirm>
          <Popconfirm
            title="确定按当前筛选条件删除所有匹配的记录吗？此操作不可恢复！"
            onConfirm={deleteByCondition}
          >
            <Button type="danger">按条件删除</Button>
          </Popconfirm>
        </Space>
      </div>

      {/* 表格 */}
      <Table
        columns={columns}
        dataSource={conversations}
        loading={loading}
        rowKey="id"
        rowSelection={{
          selectedRowKeys: selectedKeys,
          onChange: setSelectedKeys,
        }}
        pagination={{
          currentPage: page,
          pageSize: pageSize,
          total: total,
          onPageChange: setPage,
          showSizeChanger: true,
          onPageSizeChange: (size) => {
            setPageSize(size);
            setPage(1);
          },
        }}
        scroll={{ x: 1400 }}
      />

      {/* 详情对话框 */}
      <Modal
        title="对话详情"
        visible={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={null}
        width={800}
        bodyStyle={{ maxHeight: '70vh', overflow: 'auto' }}
      >
        {currentDetail && (
          <div>
            <h4>基本信息</h4>
            <p>
              <strong>用户:</strong> {currentDetail.username} (ID: {currentDetail.user_id})
            </p>
            <p>
              <strong>模型:</strong> {currentDetail.model_name}
            </p>
            <p>
              <strong>Token使用:</strong> 输入 {currentDetail.prompt_tokens} / 输出{' '}
              {currentDetail.completion_tokens} / 总计 {currentDetail.total_tokens}
            </p>
            <p>
              <strong>响应时间:</strong> {currentDetail.use_time} ms
            </p>
            <p>
              <strong>创建时间:</strong> {renderTimestamp(currentDetail.created_at)}
            </p>

            <h4 style={{ marginTop: 20 }}>请求消息</h4>
            <pre
              style={{
                background: '#f6f7f9',
                padding: 12,
                borderRadius: 4,
                maxHeight: 300,
                overflow: 'auto',
              }}
            >
              {JSON.stringify(JSON.parse(currentDetail.request_messages || '[]'), null, 2)}
            </pre>

            <h4 style={{ marginTop: 20 }}>响应内容</h4>
            <pre
              style={{
                background: '#f6f7f9',
                padding: 12,
                borderRadius: 4,
                maxHeight: 300,
                overflow: 'auto',
              }}
            >
              {currentDetail.response_content}
            </pre>
          </div>
        )}
      </Modal>

      {/* 设置对话框 */}
      <Modal
        title="对话记录功能设置"
        visible={settingVisible}
        onCancel={() => setSettingVisible(false)}
        onOk={() => updateSetting(!conversationLogEnabled)}
        okText="确定"
        cancelText="取消"
      >
        <p>当前状态：{conversationLogEnabled ? '已启用' : '已禁用'}</p>
        <p>是否要{conversationLogEnabled ? '禁用' : '启用'}对话记录功能？</p>
        {!conversationLogEnabled && (
          <p style={{ color: '#f5222d' }}>
            警告：启用后将记录所有用户的对话内容，请确保符合隐私政策和法律法规。
          </p>
        )}
      </Modal>
    </div>
  );
};

export default ConversationManagement;
