/**
 * 安全工具函数
 * 用于处理敏感数据的显示和存储
 */

/**
 * 掩码 API 密钥，只显示前后几位
 * @param key - 完整的 API 密钥
 * @param prefixLength - 显示前几位（默认 4）
 * @param suffixLength - 显示后几位（默认 4）
 * @returns 掩码后的密钥，如 "sk-ab***...***xyz"
 */
export function maskApiKey(key: string, prefixLength = 4, suffixLength = 4): string {
  if (!key || key.length <= prefixLength + suffixLength) {
    return '***';
  }

  const prefix = key.substring(0, prefixLength);
  const suffix = key.substring(key.length - suffixLength);
  
  return `${prefix}***...***${suffix}`;
}

/**
 * 验证 API 密钥格式
 * @param key - API 密钥
 * @param provider - 提供商类型
 * @returns 是否有效
 */
export function validateApiKey(key: string, provider: string): boolean {
  if (!key || key.trim().length === 0) {
    return false;
  }

  switch (provider) {
    case 'openai':
      // OpenAI 密钥格式: sk-...（至少 20 个字符）
      return key.startsWith('sk-') && key.length >= 20;
    
    case 'deepseek':
      // DeepSeek 密钥格式可能不同，这里做基本验证
      return key.length >= 20;
    
    case 'custom':
      // 自定义提供商，只检查长度
      return key.length >= 10;
    
    default:
      return key.length >= 10;
  }
}

/**
 * 清理输入，防止 XSS 攻击
 * @param input - 用户输入
 * @returns 清理后的字符串
 */
export function sanitizeInput(input: string): string {
  if (!input) return '';
  
  return input
    .trim()
    // 移除 HTML 标签
    .replace(/<[^>]*>/g, '')
    // 移除脚本标签内容
    .replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')
    // 转义特殊字符
    .replace(/[<>'"]/g, (char) => {
      const escapeMap: Record<string, string> = {
        '<': '&lt;',
        '>': '&gt;',
        "'": '&#39;',
        '"': '&quot;',
      };
      return escapeMap[char] || char;
    });
}

/**
 * 验证字符串长度
 * @param value - 输入值
 * @param min - 最小长度
 * @param max - 最大长度
 * @returns 是否有效
 */
export function validateLength(value: string, min: number, max: number): boolean {
  const length = value.trim().length;
  return length >= min && length <= max;
}

/**
 * 验证数字范围
 * @param value - 输入值
 * @param min - 最小值
 * @param max - 最大值
 * @returns 是否有效
 */
export function validateNumber(value: string | number, min: number, max: number): boolean {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  return !isNaN(num) && num >= min && num <= max;
}

/**
 * 验证食物名称
 * @param name - 食物名称
 * @returns 是否有效
 */
export function validateFoodName(name: string): boolean {
  // 2-100 个字符，只允许字母、数字、空格、连字符和中文
  return (
    validateLength(name, 2, 100) &&
    /^[\w\s\-\u4e00-\u9fa5]+$/.test(name)
  );
}

/**
 * 验证营养值
 * @param value - 营养值
 * @returns 是否有效
 */
export function validateNutritionValue(value: string | number): boolean {
  return validateNumber(value, 0, 10000);
}

/**
 * 验证价格
 * @param value - 价格
 * @returns 是否有效
 */
export function validatePrice(value: string | number): boolean {
  return validateNumber(value, 0, 100000);
}

/**
 * 生成随机字符串（用于 CSRF token 等）
 * @param length - 长度
 * @returns 随机字符串
 */
export function generateRandomString(length: number = 32): string {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';
  const randomValues = new Uint8Array(length);
  
  if (typeof window !== 'undefined' && window.crypto) {
    window.crypto.getRandomValues(randomValues);
    for (let i = 0; i < length; i++) {
      result += chars[randomValues[i] % chars.length];
    }
  } else {
    // Fallback for server-side
    for (let i = 0; i < length; i++) {
      result += chars[Math.floor(Math.random() * chars.length)];
    }
  }
  
  return result;
}

/**
 * 检查是否为安全的 URL
 * @param url - URL 字符串
 * @returns 是否安全
 */
export function isSafeUrl(url: string): boolean {
  try {
    const parsed = new URL(url);
    // 只允许 http 和 https 协议
    return ['http:', 'https:'].includes(parsed.protocol);
  } catch {
    return false;
  }
}

/**
 * 密码强度验证
 * @param password - 密码
 * @returns 强度等级和提示
 */
export function validatePasswordStrength(password: string): {
  strength: 'weak' | 'medium' | 'strong';
  message: string;
  isValid: boolean;
} {
  if (password.length < 8) {
    return {
      strength: 'weak',
      message: 'Password must be at least 8 characters',
      isValid: false,
    };
  }

  let score = 0;
  
  // 检查长度
  if (password.length >= 12) score++;
  
  // 检查是否包含小写字母
  if (/[a-z]/.test(password)) score++;
  
  // 检查是否包含大写字母
  if (/[A-Z]/.test(password)) score++;
  
  // 检查是否包含数字
  if (/\d/.test(password)) score++;
  
  // 检查是否包含特殊字符
  if (/[!@#$%^&*(),.?":{}|<>]/.test(password)) score++;

  if (score <= 2) {
    return {
      strength: 'weak',
      message: 'Weak password. Add uppercase, numbers, and special characters.',
      isValid: true,
    };
  } else if (score <= 3) {
    return {
      strength: 'medium',
      message: 'Medium strength password.',
      isValid: true,
    };
  } else {
    return {
      strength: 'strong',
      message: 'Strong password!',
      isValid: true,
    };
  }
}

/**
 * 防止 XSS 的 HTML 编码
 * @param text - 文本
 * @returns 编码后的文本
 */
export function encodeHtml(text: string): string {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

/**
 * 安全的 JSON 解析
 * @param json - JSON 字符串
 * @param fallback - 解析失败时的默认值
 * @returns 解析结果
 */
export function safeJsonParse<T>(json: string, fallback: T): T {
  try {
    return JSON.parse(json) as T;
  } catch {
    return fallback;
  }
}
