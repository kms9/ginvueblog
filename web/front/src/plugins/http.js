import Vue from 'vue'
import axios from 'axios'

// axios请求地址

axios.defaults.baseURL = 'http://qqt.feibooks.com/api/v1'
Vue.prototype.$http = axios
