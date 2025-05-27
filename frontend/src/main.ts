import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import '@/assets/main.css'

// PrimeVue imports
import PrimeVue from 'primevue/config'
import { definePreset } from '@primeuix/themes'
import Aura from '@primeuix/themes/aura'
import 'primeicons/primeicons.css'

// PrimeVue Components
import Button from 'primevue/button'
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import DataTable from 'primevue/datatable'
import Checkbox from 'primevue/checkbox'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import Dropdown from 'primevue/select'
import Textarea from 'primevue/textarea'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Toast from 'primevue/toast'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import ConfirmDialog from 'primevue/confirmdialog'
import Menubar from 'primevue/menubar'
import Avatar from 'primevue/avatar'
import Divider from 'primevue/divider'
import Message from 'primevue/message'
import Tabs from 'primevue/tabs'
import TabList from 'primevue/tablist'
import Tab from 'primevue/tab'
import TabPanels from 'primevue/tabpanels'
import TabPanel from 'primevue/tabpanel'
import Tooltip from 'primevue/tooltip'
import Chip from 'primevue/chip'
import Timeline from 'primevue/timeline'
import MultiSelect from 'primevue/multiselect'
import Fieldset from 'primevue/fieldset'
import Panel from 'primevue/panel'
import ProgressSpinner from 'primevue/progressspinner'
import Breadcrumb from 'primevue/breadcrumb'
import Badge from 'primevue/badge'
import Chart from 'primevue/chart'

const app = createApp(App)

app.use(createPinia())
app.use(router)
const MyPreset = definePreset(Aura, {
  semantic: {
    primary: {
      50: '{blue.50}',
      100: '{blue.100}',
      200: '{blue.200}',
      300: '{blue.300}',
      400: '{blue.400}',
      500: '{blue.500}',
      600: '{blue.600}',
      700: '{blue.700}',
      800: '{blue.800}',
      900: '{blue.900}',
      950: '{blue.950}'
    }
  }
})

app.use(PrimeVue, {
  theme: {
    preset: MyPreset,
    // options: {
    //   darkModeSelector: '.dark',
    //   cssLayer: {
    //     name: 'primevue',
    //     order: 'tailwind-base, primevue, tailwind-utilities'
    //   }
    // }
  }
})
app.use(ToastService)
app.use(ConfirmationService)

// Register PrimeVue components globally
app.component('Button', Button)
app.component('Card', Card)
app.component('InputText', InputText)
app.component('Password', Password)
app.component('DataTable', DataTable)
app.component('Column', Column)
app.component('Checkbox', Checkbox)
app.component('Dialog', Dialog)
app.component('Dropdown', Dropdown)
app.component('Textarea', Textarea)
app.component('InputNumber', InputNumber)
app.component('Tag', Tag)
app.component('Toast', Toast)
app.component('ConfirmDialog', ConfirmDialog)
app.component('Menubar', Menubar)
app.component('Avatar', Avatar)
app.component('Divider', Divider)
app.component('Message', Message)
app.component('Tabs', Tabs)
app.component('TabList', TabList)
app.component('Tab', Tab)
app.component('TabPanels', TabPanels)
app.component('TabPanel', TabPanel)
app.directive('Tooltip', Tooltip)
app.component('Chip', Chip)
app.component('Timeline', Timeline)
app.component('MultiSelect', MultiSelect)
app.component('Fieldset', Fieldset)
app.component('Panel', Panel)
app.component('ProgressSpinner', ProgressSpinner)
app.component('Breadcrumb', Breadcrumb)
app.component('Badge', Badge)
app.component('Chart', Chart)

app.mount('#app')
