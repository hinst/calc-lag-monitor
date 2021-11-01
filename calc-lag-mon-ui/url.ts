export function buildParametersString(params: { [key: string]: string }) {
  let text = '';
  for (const key in params)
    text += (text.length ? '&' : '?') +
      encodeURIComponent(key) + '=' + encodeURIComponent(params[key]);
  return text;
}
