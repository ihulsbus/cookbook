import { createOidcAuth, SignInType, LogLevel } from 'vue-oidc-client/vue3';

const loco = window.location;
const appRootUrl = `${loco.protocol}//${loco.host}${process.env.BASE_URL}`;

const idsrvAuth = createOidcAuth(
  'main',
  SignInType.Window,
  appRootUrl,
  {
    authority: process.env.VUE_APP_COOKBOOK_AUTHORITY,
    client_id: process.env.VUE_APP_COOKBOOK_CLIENTID,
    response_type: 'code',
    redirect_uri: `${appRootUrl}oidc/callback`,
    scope: 'openid profile email',
    prompt: 'login',
  },
  console,
  LogLevel.Info,
);

export default idsrvAuth;
