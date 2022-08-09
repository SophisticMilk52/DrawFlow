import {createStore} from "vuex";
Vue.use(Vuex)
const store = createStore({
    state: {
        sender: [],
        subscriber: []
    },
    getters: {
        getSender(state) {
            return state.sender
        },
        getSubscriber(state) {
            return state.subscriber
        }
    }
});

export default store;