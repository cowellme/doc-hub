import { render, screen } from '@testing-library/react';
import App from './App';

test('renders upload heading', () => {
  render(<App />);
  const headingElement = screen.getByText(/upload a pdf or docx/i);
  expect(headingElement).toBeInTheDocument();
});
