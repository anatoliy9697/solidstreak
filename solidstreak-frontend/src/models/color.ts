export interface Color {
    name: string
    value50: string
    value50hex: string
    value100: string
    value100hex: string
    value200: string
    value200hex: string
    value400: string
    value400hex: string
    value500: string
    value500hex: string
    value600: string
    value600hex: string
    value800: string
    value800hex: string
}


export const RED: Color = {
    name: 'red',
    value50: 'red-50',
    value50hex: '#fef2f2',
    value100: 'red-100',
    value100hex: '#fee2e2',
    value200: 'red-200',
    value200hex: '#fecaca',
    value400: 'red-400',
    value400hex: '#f87171',
    value500: 'red-500',
    value500hex: '#ef4444',
    value600: 'red-600',
    value600hex: '#dc2626',
    value800: 'red-800',
    value800hex: '#991b1b',
};
export const ORANGE: Color = {
    name: 'orange',
    value50: 'orange-50',
    value50hex: '#fff7ed',
    value100: 'orange-100',
    value100hex: '#ffedd5',
    value200: 'orange-200',
    value200hex: '#fed7aa',
    value400: 'orange-400',
    value400hex: '#fb923c',
    value500: 'orange-500',
    value500hex: '#f97316',
    value600: 'orange-600',
    value600hex: '#ea580c',
    value800: 'orange-800',
    value800hex: '#9a3412',
};
export const YELLOW: Color = {
    name: 'yellow',
    value50: 'yellow-50',
    value50hex: '#fefce8',
    value100: 'yellow-100',
    value100hex: '#fef9c3',
    value200: 'yellow-200',
    value200hex: '#fef08a',
    value400: 'yellow-400',
    value400hex: '#facc15',
    value500: 'yellow-500',
    value500hex: '#eab308',
    value600: 'yellow-600',
    value600hex: '#ca8a04',
    value800: 'yellow-800',
    value800hex: '#854d0e',
};
export const LIME: Color = {
    name: 'lime',
    value50: 'lime-50',
    value50hex: '#f7fee7',
    value100: 'lime-100',
    value100hex: '#ecfccb',
    value200: 'lime-200',
    value200hex: '#d9f99d',
    value400: 'lime-400',
    value400hex: '#a3e635',
    value500: 'lime-500',
    value500hex: '#84cc16',
    value600: 'lime-600',
    value600hex: '#65a30d',
    value800: 'lime-800',
    value800hex: '#365314',
};
export const GREEN: Color = {
    name: 'green',
    value50: 'green-50',
    value50hex: '#f0fdf4',
    value100: 'green-100',
    value100hex: '#dcfce7',
    value200: 'green-200',
    value200hex: '#bbf7d0',
    value400: 'green-400',
    value400hex: '#4ade80',
    value500: 'green-500',
    value500hex: '#22c55e',
    value600: 'green-600',
    value600hex: '#16a34a',
    value800: 'green-800',
    value800hex: '#166534',
};
export const BLUE: Color = {
    name: 'blue',
    value50: 'blue-50',
    value50hex: '#eff6ff',
    value100: 'blue-100',
    value100hex: '#dbeafe',
    value200: 'blue-200',
    value200hex: '#bfdbfe',
    value400: 'blue-400',
    value400hex: '#60a5fa',
    value500: 'blue-500',
    value500hex: '#3b82f6',
    value600: 'blue-600',
    value600hex: '#2563eb',
    value800: 'blue-800',
    value800hex: '#1e40af',
};
export const PURPLE: Color = {
    name: 'purple',
    value50: 'purple-50',
    value50hex: '#faf5ff',
    value100: 'purple-100',
    value100hex: '#f3e8ff',
    value200: 'purple-200',
    value200hex: '#e9d5ff',
    value400: 'purple-400',
    value400hex: '#a78bfa',
    value500: 'purple-500',
    value500hex: '#8b5cf6',
    value600: 'purple-600',
    value600hex: '#7c3aed',
    value800: 'purple-800',
    value800hex: '#6d28d9',
};

export const COLORS: { [key: string]: Color } = {
    [RED.name]: RED,
    [ORANGE.name]: ORANGE,
    [YELLOW.name]: YELLOW,
    [LIME.name]: LIME,
    [GREEN.name]: GREEN,
    [BLUE.name]: BLUE,
    [PURPLE.name]: PURPLE
};

function hexToRgb(hex: string) {
    const cleanHex = hex.replace('#', '');
    const bigint = parseInt(cleanHex, 16);
    return {
        r: (bigint >> 16) & 255,
        g: (bigint >> 8) & 255,
        b: bigint & 255,
    };
};

function rgbToHex(r: number, g: number, b: number) {
    const toHex = (c: number) => c.toString(16).padStart(2, '0');
    return `#${toHex(r)}${toHex(g)}${toHex(b)}`;
};

export function generateColorGradient(leftHex: string, rightHex: string, n: number): string[] {
    if (n <= 0) {
        return [];
    } else if (n === 1) {
        return [rightHex];
    } else if (n === 2) {
        return [leftHex, rightHex];
    }

    const left = hexToRgb(leftHex);
    const right = hexToRgb(rightHex);

    const result: string[] = [];

    for (let i = 0; i < n; i++) {
        const t = i / (n - 1);
        const r = Math.round(left.r + (right.r - left.r) * t);
        const g = Math.round(left.g + (right.g - left.g) * t);
        const b = Math.round(left.b + (right.b - left.b) * t);
        result.push(rgbToHex(r, g, b));
    }

    return result;
}