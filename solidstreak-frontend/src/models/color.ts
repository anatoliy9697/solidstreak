export interface Color {
    value50: string
    value50hex: string
    value100: string
    value100hex: string
    value400: string
    value400hex: string
    value500: string
    value500hex: string
    value600: string
    value600hex: string
}

export const GREEN: Color = {
    value50: 'green-50',
    value50hex: '#f0fdf4',
    value100: 'green-100',
    value100hex: '#dcfce7',
    value400: 'green-400',
    value400hex: '#4ade80',
    value500: 'green-500',
    value500hex: '#22c55e',
    value600: 'green-600',
    value600hex: '#16a34a',
};
export const RED: Color = {
    value50: 'red-50',
    value50hex: '#fef2f2',
    value100: 'red-100',
    value100hex: '#fee2e2',
    value400: 'red-400',
    value400hex: '#f87171',
    value500: 'red-500',
    value500hex: '#ef4444',
    value600: 'red-600',
    value600hex: '#dc2626',
};
export const BLUE: Color = {
    value50: 'blue-50',
    value50hex: '#eff6ff',
    value100: 'blue-100',
    value100hex: '#dbeafe',
    value400: 'blue-400',
    value400hex: '#60a5fa',
    value500: 'blue-500',
    value500hex: '#3b82f6',
    value600: 'blue-600',
    value600hex: '#2563eb',
};
export const YELLOW: Color = {
    value50: 'yellow-50',
    value50hex: '#fefce8',
    value100: 'yellow-100',
    value100hex: '#fef9c3',
    value400: 'yellow-400',
    value400hex: '#facc15',
    value500: 'yellow-500',
    value500hex: '#eab308',
    value600: 'yellow-600',
    value600hex: '#ca8a04',
};

export const COLORS: { [key: string]: Color } = {
    'green': GREEN,
    'red': RED,
    'blue': BLUE,
    'yellow': YELLOW,
};