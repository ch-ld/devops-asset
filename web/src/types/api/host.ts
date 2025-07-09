// Types for CMDB Provider and Host

export type ProviderType = 'aliyun' | 'tencent' | 'huawei' | 'aws' | 'volcengine';

export interface Provider {
  id: number;
  name: string;
  type: ProviderType;
  access_key_id: string;
  status: boolean;
  remark: string;
  created_at: number;
  updated_at: number;
}

export interface Host {
  id: number;
  provider_id: number;
  instance_id: string;
  name: string;
  resource_type: string;
  region: string;
  public_ip: string[];
  private_ip: string[];
  configuration: {
    cpu: number;
    memory: string; // e.g., "4GB"
    disk: string;   // e.g., "100GB"
  };
  os: string;
  status: string;
  expired_at: number | null;
  remark: string;
  created_at: number;
  updated_at: number;
} 