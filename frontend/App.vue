<template>
  <div class="container">
    <div v-if="currentEpisode">
      <p class="lead">{{ currentEpisode.title }}</p>
      <p>{{ currentEpisode.description }}</p>
      <audio controls id="audio" :src="currentEpisode.url" controlslist="nodownload" style="width: 100%;"></audio>
    </div>
    <hr>
    <div class="card">
      <div class="card-header">
        サブスクリプション
        <div class="float-right">
          {{ crawlState }}
          <button class="btn btn-secondary btn-sm" @click="crawlSubscriptions"
            :disabled="crawlSubscriptionState">すべて更新する</button>
          <a class="btn btn-primary btn-sm" @click="fetchSubscriptions">再読み込み</a>
        </div>
      </div>
      <ul class="list-group">
        <li v-for="(item, index) in subscriptions" :key="index" class="list-group-item">
          <a class="link-primary" style="cursor: pointer" :data-index="index" @click="openSubscription">{{ item.title
            }}</a>
          <span class="badge badge-light" v-if="item.new_entry_count > 0">
            {{ item.new_entry_count }}
          </span>
          <div class="float-right">
            <a class="btn btn-sm btn-warning" @click="deleteSubscriptionPre" :data-index="index">削除</a>
          </div>
        </li>
      </ul>
    </div>

    <div v-if="episodes.length > 0">
      <hr>
      <div class="card">
        <div class="card-header">
          エピソード
          <div class="float-right">
            <a class="btn btn-secondary btn-sm" @click="fetchSubscription">更新</a>
          </div>
        </div>
        <ul class="list-group">
          <li v-for="(item, index) in episodes" :key="index" class="list-group-item">
            <span class="badge badge-light" v-if="item.new">NEW</span>
            <a class="link-primary" style="cursor: pointer" :data-index="index" @click="openEpisode">{{ item.title
              }}</a>
            <div class="float-right">
              <a class="btn btn-sm btn-warning" @click="deleteEpisode" :data-index="index">削除</a>
            </div>
          </li>
        </ul>
      </div>
    </div>

    <div>
      <hr>
      <form>
        <div class="form-group">
          <label for="subUrl">Podcast Feed URL</label>
          <input type="url" class="form-control" v-model="feedURL" aria-describedby="subscription URL">
          <a @click="registerSubscription" class="btn btn-primary" v-if="feedURL.length > 0">登録する</a>
        </div>
      </form>
    </div>

    <div class="modal" tabindex="-1" style="display:block;" v-if="!!deleteSubscriptionVal">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">サブスクリプションの削除</h5>
          </div>
          <div class="modal-body">
            <p>{{ deleteSubscriptionVal.title }} を削除しますか?</p>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-dismiss="modal"
              @click="deleteSubscriptionVal = undefined">キャンセル</button>
            <button type="button" class="btn btn-warning" @click="deleteSubscription">削除する</button>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue"

type Subscription = {
  id: number
  title: string
  url: string
  new_entry_count: number
}

type Episode = {
  id: number
  title: string
  description: string
  url: string
  new: boolean
}

import createClient from "openapi-fetch"
import type { paths } from "./api/schema"

const client = createClient<paths>({})

export default defineComponent({
  setup: () => {
    const subscriptions = ref<Subscription[]>([])

    const currentSubscription = ref<number | undefined>(undefined)
    const episodes = ref<Episode[]>([])

    const currentEpisode = ref<Episode | undefined>(undefined)

    const feedURL = ref<string>("")

    const crawlState = ref<string>("")
    const crawlSubscriptionState = ref<boolean>(false)

    const deleteSubscriptionVal = ref<Subscription | undefined>(undefined)

    const registerSubscription = () => {
      client.POST("/api/subscription", {
        body: {
          url: feedURL.value,
        }
      }).then((data) => {
        console.log(data)
        fetchSubscriptions()
        feedURL.value = ""
      }).catch((error) => {
        alert(error)
        console.log(error)
      })
    }

    const readSubscription = (id: number) => {
      client.GET("/api/subscription/{id}", {
        params: {
          path: {
            id
          }
        },
      }).then((data) => {
        if (data.data === undefined) {
          return
        }
        episodes.value = data?.data?.sort((a, b) => (b.publishedAt < a.publishedAt ? -1 : 1))
        currentSubscription.value = id
      })
    }

    const fetchSubscription = () => {
      const id = currentSubscription.value
      if (id === undefined) {
        return
      }

      client.POST("/api/subscription/{id}/-/fetch", {
        params: {
          path: {
            id
          }
        },
      }).then((data) => {
        console.log(data)
        readSubscription(id)
      })
    }

    const openSubscription = (event: Event) => {
      const target = event.target as HTMLElement

      const dataIndex = target.getAttribute("data-index")
      if (dataIndex === null) {
        return
      }
      const index = parseInt(dataIndex, 10)
      const subscription = subscriptions.value[index]

      if (subscription === undefined) {
        return
      }

      readSubscription(subscription.id)
    }

    const openEpisode = (event: Event) => {
      const target = event.target as HTMLElement

      const dataIndex = target.getAttribute("data-index")
      if (dataIndex === null) {
        return
      }
      const index = parseInt(dataIndex, 10)
      const episode = episodes.value[index]
      if (episode === undefined) {
        return
      }

      const subId = currentSubscription.value
      if (subId === undefined) {
        return
      }

      client.POST("/api/subscription/{id}/{entryId}/open", {
        params: {
          path: {
            id: subId,
            entryId: episode.id,
          }
        },
      }).then((data) => {
        console.log(data.data)
      })

      currentEpisode.value = episode
    }

    const deleteSubscriptionPre = (event: Event) => {
      const target = event.target as HTMLElement

      const dataIndex = target.getAttribute("data-index")
      if (dataIndex === null) {
        return
      }
      const index = parseInt(dataIndex, 10)

      const subscription = subscriptions.value[index]
      if (subscription === undefined) {
        return
      }

      deleteSubscriptionVal.value = subscription
    }

    const deleteSubscription = () => {
      const v = deleteSubscriptionVal.value
      if (v === undefined) {
        return
      }
      client.DELETE("/api/subscription/{id}", {
        params: {
          path: {
            id: v.id,
          }
        },
      }).then((data) => {
        console.log(data.data)
        deleteSubscriptionVal.value = undefined
        fetchSubscriptions()
      })
    }

    const deleteEpisode = (event: Event) => {
      const target = event.target as HTMLElement

      const dataIndex = target.getAttribute("data-index")
      if (dataIndex === null) {
        return
      }
      const index = parseInt(dataIndex, 10)
      if (index === 0) {
        alert("先頭のエピソードは削除できません")
        return
      }

      const episode = episodes.value[index]
      if (episode === undefined) {
        return
      }

      const subId = currentSubscription.value
      if (subId === undefined) {
        return
      }

      client.DELETE("/api/subscription/{id}/{entryId}", {
        params: {
          path: {
            id: subId,
            entryId: episode.id,
          }
        },
      }).then((data) => {
        console.log(data.data)
      })
      episodes.value.splice(index, 1)
    }

    const fetchSubscriptions = () => {
      client.GET("/api/subscriptions", {}).then((data) => {
        if (data.data === undefined) {
          return
        }
        subscriptions.value = data.data
      })
    }

    const crawlSubscriptions = async () => {
      crawlSubscriptionState.value = true
      crawlState.value = "0" + "/" + subscriptions.value.length

      const f = async (id: number, index: number) => {
        return new Promise<void>(async (resolve) => {
          await client.POST("/api/subscription/{id}/-/fetch", {
            params: {
              path: {
                id
              }
            }
          }).then(() => {
            crawlState.value = (index + 1) + "/" + subscriptions.value.length
          })
          resolve()
        })
      }

      for (const [index, sub] of subscriptions.value.entries()) await f(sub.id, index)

      crawlState.value = ""
      crawlSubscriptionState.value = false

      fetchSubscriptions()
    }

    // on load
    fetchSubscriptions()

    return {
      feedURL,

      crawlState,
      crawlSubscriptionState,

      deleteSubscriptionVal,

      subscriptions,
      episodes,
      currentEpisode,

      openSubscription,
      deleteSubscriptionPre,
      deleteSubscription,

      openEpisode,
      deleteEpisode,

      fetchSubscription,

      registerSubscription,

      fetchSubscriptions,
      crawlSubscriptions,
    }
  },
})
</script>
