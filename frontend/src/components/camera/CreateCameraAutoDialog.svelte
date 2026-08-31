<script lang="ts">
  import { models } from "../../../wailsjs/go/models";
  import { GetRTSPStreamURL, AddCamera } from "../../../wailsjs/go/main/App";
  let { selectedCamera = $bindable<models.DiscoveredCamera | null>(null) } =
    $props();
  let username: string = $state("");
  let password: string = $state("");

  async function CreateACamera() {
    try {
      let streaming_url = await GetRTSPStreamURL(
        selectedCamera.endpoint,
        username,
        password,
      );
      console.log("CACACAufiziauaza", selectedCamera.endpoint);
      await AddCamera(selectedCamera.name, streaming_url);
    } catch (e) {
      console.error("Erreur lors de la récupération du stream:", e);
    }
  }
</script>

<button
  class="btn-back"
  onclick={() => {
    selectedCamera = null;
  }}>&larr; Retour</button
>
<h2>Ajouter la caméra</h2>

<form
  onsubmit={(e) => {
    e.preventDefault();
    CreateACamera();
  }}
>
  <div class="form-group">
    <label for="camera-name">Nom de la caméra</label>
    <input type="text" id="camera-name" value={selectedCamera?.name || ""} />
  </div>

  <div class="form-group">
    <label for="camera-endpoint">Endpoint</label>
    <input
      type="text"
      id="camera-endpoint"
      value={selectedCamera?.endpoint || ""}
    />
  </div>

  <div class="form-group">
    <label for="camera-username">Utilisateur</label>
    <input
      type="text"
      id="camera-username"
      bind:value={username}
      autocomplete="username"
    />
  </div>
  <div class="form-group">
    <label for="camera-password">Mot de passe</label>
    <input
      type="password"
      id="camera-password"
      bind:value={password}
      autocomplete="current-password"
    />
  </div>
  <button type="submit" class="btn-submit">Créer</button>
</form>

<style>
  .btn-back {
    all: unset;
    font-size: 0.8rem;
    color: #888;
    cursor: pointer;
    margin-bottom: 0.75rem;
  }

  .btn-back:hover {
    color: #222;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
    margin-bottom: 1.1rem;
  }

  .form-group label {
    font-size: 0.8rem;
    color: #666;
  }

  .form-group input {
    box-sizing: border-box;
    font-size: 0.9rem;
    padding: 0.6rem 0.7rem;
    border: 1px solid #ddd;
    border-radius: var(--radius-sm);
  }

  .form-group input:focus {
    outline: none;
    border-color: #999;
  }

  .btn-submit {
    all: unset;
    box-sizing: border-box;
    background: #1a1a1a;
    color: #fff;
    font-size: 0.85rem;
    font-weight: 500;
    padding: 0.6rem 1rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: background-color 0.15s ease;
  }

  .btn-submit:hover {
    background: #333;
  }
</style>
