import { useState, type KeyboardEvent } from 'react';

interface TagInputProps {
    value: string[];
    onChange: (tags: string[]) => void;
    placeholder?: string;
}

export default function TagInput({
    value,
    onChange,
    placeholder = '입력 후 Enter',
}: TagInputProps) {
    const [input, setInput] = useState('');

    const add = () => {
        const trimmed = input.trim();
        if (trimmed && !value.includes(trimmed)) {
            onChange([...value, trimmed]);
        }
        setInput('');
    };

    const remove = (tag: string) => onChange(value.filter((t) => t !== tag));

    const handleKey = (e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter' || e.key === ',') {
            e.preventDefault();
            add();
        } else if (e.key === 'Backspace' && !input && value.length > 0) {
            onChange(value.slice(0, -1));
        }
    };

    return (
        <div className="flex flex-wrap gap-1.5 min-h-[36px] px-2 py-1.5 bg-surface border border-border rounded-lg focus-within:ring-2 focus-within:ring-primary/30 focus-within:border-primary transition-colors">
            {value.map((tag) => (
                <span
                    key={tag}
                    className="inline-flex items-center gap-1 px-2 py-0.5 bg-tag text-tag-text text-[11px] font-medium rounded-full"
                >
                    {tag}
                    <button
                        type="button"
                        onClick={() => remove(tag)}
                        className="hover:text-primary-dark leading-none"
                    >
                        ×
                    </button>
                </span>
            ))}
            <input
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={handleKey}
                onBlur={add}
                placeholder={value.length === 0 ? placeholder : ''}
                className="flex-1 min-w-[80px] text-sm text-text bg-transparent outline-none placeholder:text-text-hint"
            />
        </div>
    );
}
