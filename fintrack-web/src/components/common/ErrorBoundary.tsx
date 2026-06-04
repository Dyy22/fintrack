import {
  Component,
  type ErrorInfo,
  type PropsWithChildren,
  type ReactNode,
} from "react";
import { Button } from "./Button";

type ErrorBoundaryProps = PropsWithChildren<{
  fallback?: ReactNode;
}>;

type ErrorBoundaryState = {
  hasError: boolean;
  error: Error | null;
};

export class ErrorBoundary extends Component<
  ErrorBoundaryProps,
  ErrorBoundaryState
> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
    console.error("ErrorBoundary caught:", error, errorInfo);
  }

  handleRetry = () => {
    this.setState({ hasError: false, error: null });
  };

  render() {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback;
      }

      return (
        <div className="flex min-h-[200px] items-center justify-center p-8">
          <div className="text-center">
            <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-red-100">
              <span className="text-2xl font-bold text-red-600">!</span>
            </div>
            <h3 className="text-lg font-semibold text-slate-950 dark:text-slate-100">
              Something went wrong
            </h3>
            <p className="mt-1 text-sm text-slate-500 dark:text-slate-400">
              An unexpected error occurred. Please try again.
            </p>
            <p className="mt-1 text-xs text-slate-400 dark:text-slate-500">
              {this.state.error?.message}
            </p>
            <Button className="mt-4" onClick={this.handleRetry}>
              Try Again
            </Button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}
