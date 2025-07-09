export interface Provider {
  id: number;
  name: string;
  type: 'aliyun' | 'tencent' | 'aws';
  access_key_id: string;
  access_key_secret: string;
  status: boolean;
  remark: string;
  created_at: string;
  updated_at: string;
} 