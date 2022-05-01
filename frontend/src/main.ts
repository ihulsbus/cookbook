import { createApp, h } from 'vue';

import PrimeVue from 'primevue/config';
import 'primeflex/primeflex.css';
import InputText from 'primevue/inputtext';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import SplitButton from 'primevue/splitbutton';
import Textarea from 'primevue/textarea';
import Dropdown from 'primevue/dropdown';
import RadioButton from 'primevue/radiobutton';
import InputNumber from 'primevue/inputnumber';
import Toast from 'primevue/toast';
import ToastService from 'primevue/toastservice';
import ConfirmationService from 'primevue/confirmationservice';
import Tooltip from 'primevue/tooltip';
import idsrvAuth from '@/lib/auth/auth';
import router from './router';
import App from './App.vue';

idsrvAuth.startup().then((ok) => {
  if (ok) {
    const app = createApp({
      render: () => h(App),
    });

    app.config.globalProperties.$oidc = idsrvAuth;

    app.use(router);
    app.use(PrimeVue, { ripple: true });
    app.use(ToastService);
    app.use(ConfirmationService);

    app.component('InputText', InputText);
    app.component('Dialog', Dialog);
    app.component('Button', Button);
    app.component('SplitButton', SplitButton);
    app.component('Textarea', Textarea);
    app.component('Dropdown', Dropdown);
    app.component('RadioButton', RadioButton);
    app.component('InputNumber', InputNumber);
    app.component('Toast', Toast);

    app.directive('tooltip', Tooltip);

    app.mount('#app');
  } else {
    console.log('Startup was not ok');
  }
});
