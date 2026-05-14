import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';

interface PostPreviewProps {
    markdown: string;
}

export default function PostPreview({ markdown }: PostPreviewProps) {
    if (!markdown) {
        return (
            <div className="flex flex-col items-center justify-center h-full text-text-hint">
                <span className="text-4xl mb-3">📝</span>
                <p className="text-sm">「미리보기 생성」을 눌러 결과를 확인하세요.</p>
            </div>
        );
    }

    return (
        <div className="prose-preview overflow-y-auto h-full">
            <ReactMarkdown remarkPlugins={[remarkGfm]}>{markdown}</ReactMarkdown>
        </div>
    );
}
