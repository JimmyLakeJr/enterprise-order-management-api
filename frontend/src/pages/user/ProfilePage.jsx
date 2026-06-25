import { useEffect, useRef, useState } from "react";
import Badge from "../../components/common/Badge";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import ErrorMessage from "../../components/common/ErrorMessage";
import Input from "../../components/common/Input";
import Toast from "../../components/common/Toast";
import { useAuth } from "../../contexts/AuthContext";
import { resolveAssetUrl } from "../../utils/resolveAssetUrl";

const IMAGE_TYPES = new Set(["image/jpeg", "image/png", "image/webp", "image/gif", "image/avif"]);
const VIDEO_TYPES = new Set(["video/mp4", "video/webm"]);
const MAX_IMAGE_SIZE = 5 * 1024 * 1024;
const MAX_VIDEO_SIZE = 20 * 1024 * 1024;
const MAX_VIDEO_DURATION = 10;
const MIN_PASSWORD_LENGTH = 6;

function validatePhone(phone) {
  const trimmed = phone.trim();
  if (!trimmed) return "";
  if (!/^\+?[0-9\s()-]+$/.test(trimmed)) return "Số điện thoại không được chứa chữ cái.";

  const normalized = trimmed.replace(/[^\d+]/g, "");
  const digitsOnly = normalized.startsWith("+") ? normalized.slice(1) : normalized;
  if (digitsOnly.length < 9 || digitsOnly.length > 15) return "Số điện thoại phải dài từ 9 đến 15 số.";
  return "";
}

export default function ProfilePage() {
  const { user, updateProfile, changePassword, uploadAvatar, uploadProfileVideo } = useAuth();
  const videoRequestRef = useRef(0);
  const [avatarUrl, setAvatarUrl] = useState("");
  const [failedAvatarUrl, setFailedAvatarUrl] = useState("");
  const [localImagePreview, setLocalImagePreview] = useState("");
  const [localVideoPreview, setLocalVideoPreview] = useState("");
  const [selectedAvatarFile, setSelectedAvatarFile] = useState(null);
  const [selectedVideoFile, setSelectedVideoFile] = useState(null);
  const [imageError, setImageError] = useState("");
  const [videoError, setVideoError] = useState("");
  const [profileForm, setProfileForm] = useState(() => ({ name: user?.name || "", phone: user?.phone || "" }));
  const [profileError, setProfileError] = useState("");
  const [profileSuccess, setProfileSuccess] = useState("");
  const [savingProfile, setSavingProfile] = useState(false);
  const [uploadingAvatar, setUploadingAvatar] = useState(false);
  const [uploadingVideo, setUploadingVideo] = useState(false);
  const [passwordForm, setPasswordForm] = useState({
    currentPassword: "",
    newPassword: "",
    confirmPassword: "",
  });
  const [passwordError, setPasswordError] = useState("");
  const [savingPassword, setSavingPassword] = useState(false);

  useEffect(() => {
    return () => {
      if (localImagePreview) URL.revokeObjectURL(localImagePreview);
    };
  }, [localImagePreview]);

  useEffect(() => {
    return () => {
      if (localVideoPreview) URL.revokeObjectURL(localVideoPreview);
    };
  }, [localVideoPreview]);

  useEffect(() => {
    return () => {
      videoRequestRef.current += 1;
    };
  }, []);

  function clearAvatarSelection() {
    if (localImagePreview) URL.revokeObjectURL(localImagePreview);
    setLocalImagePreview("");
    setSelectedAvatarFile(null);
  }

  function clearVideoSelection() {
    if (localVideoPreview) URL.revokeObjectURL(localVideoPreview);
    setLocalVideoPreview("");
    setSelectedVideoFile(null);
  }

  function handleImageFile(event) {
    const file = event.target.files?.[0];
    setImageError("");

    if (!file) {
      clearAvatarSelection();
      return;
    }
    if (!IMAGE_TYPES.has(file.type)) {
      setImageError("Chỉ chấp nhận JPG/JPEG, PNG, WebP, GIF hoặc AVIF.");
      event.target.value = "";
      return;
    }
    if (file.size > MAX_IMAGE_SIZE) {
      setImageError("Ảnh không được vượt quá 5 MB.");
      event.target.value = "";
      return;
    }

    clearAvatarSelection();
    setSelectedAvatarFile(file);
    setLocalImagePreview(URL.createObjectURL(file));
  }

  function handleVideoFile(event) {
    const file = event.target.files?.[0];
    const requestId = ++videoRequestRef.current;
    setVideoError("");

    if (!file) {
      clearVideoSelection();
      return;
    }
    if (!VIDEO_TYPES.has(file.type)) {
      setVideoError("Chỉ chấp nhận video MP4 hoặc WebM.");
      event.target.value = "";
      return;
    }
    if (file.size > MAX_VIDEO_SIZE) {
      setVideoError("Video không được vượt quá 20 MB.");
      event.target.value = "";
      return;
    }

    const objectUrl = URL.createObjectURL(file);
    const probe = document.createElement("video");
    probe.preload = "metadata";
    probe.onloadedmetadata = () => {
      if (requestId !== videoRequestRef.current) {
        URL.revokeObjectURL(objectUrl);
        return;
      }
      if (!Number.isFinite(probe.duration) || probe.duration > MAX_VIDEO_DURATION) {
        setVideoError("Video phải có thời lượng tối đa 10 giây.");
        URL.revokeObjectURL(objectUrl);
        event.target.value = "";
        return;
      }

      clearVideoSelection();
      setSelectedVideoFile(file);
      setLocalVideoPreview(objectUrl);
    };
    probe.onerror = () => {
      if (requestId === videoRequestRef.current) {
        setVideoError("Không đọc được thông tin video.");
        event.target.value = "";
      }
      URL.revokeObjectURL(objectUrl);
    };
    probe.src = objectUrl;
  }

  async function handleProfileSubmit(event) {
    event.preventDefault();

    const name = profileForm.name.trim();
    const phoneValidationError = validatePhone(profileForm.phone);

    if (name.length < 2) {
      setProfileError("Tên hiển thị phải có ít nhất 2 ký tự.");
      return;
    }
    if (phoneValidationError) {
      setProfileError(phoneValidationError);
      return;
    }

    setSavingProfile(true);
    setProfileError("");
    try {
      const updatedUser = await updateProfile({
        name,
        phone: profileForm.phone.trim(),
      });
      setProfileForm({
        name: updatedUser?.name || "",
        phone: updatedUser?.phone || "",
      });
      setProfileSuccess("Đã cập nhật hồ sơ.");
    } catch (err) {
      setProfileError(err?.response?.data?.message || "Không thể cập nhật hồ sơ.");
    } finally {
      setSavingProfile(false);
    }
  }

  async function handlePasswordSubmit(event) {
    event.preventDefault();

    const currentPassword = passwordForm.currentPassword;
    const newPassword = passwordForm.newPassword;
    const confirmPassword = passwordForm.confirmPassword;

    if (!currentPassword || !newPassword || !confirmPassword) {
      setPasswordError("Vui lòng nhập đầy đủ thông tin đổi mật khẩu.");
      return;
    }
    if (newPassword.length < MIN_PASSWORD_LENGTH) {
      setPasswordError(`Mật khẩu mới phải có ít nhất ${MIN_PASSWORD_LENGTH} ký tự.`);
      return;
    }
    if (currentPassword === newPassword) {
      setPasswordError("Mật khẩu mới phải khác mật khẩu hiện tại.");
      return;
    }
    if (newPassword !== confirmPassword) {
      setPasswordError("Mật khẩu xác nhận không khớp.");
      return;
    }

    setSavingPassword(true);
    setPasswordError("");
    try {
      await changePassword({
        current_password: currentPassword,
        new_password: newPassword,
      });
      setPasswordForm({
        currentPassword: "",
        newPassword: "",
        confirmPassword: "",
      });
      setProfileSuccess("Đổi mật khẩu thành công.");
    } catch (err) {
      setPasswordError(err?.response?.data?.message || "Không thể đổi mật khẩu.");
    } finally {
      setSavingPassword(false);
    }
  }

  async function handleAvatarUpload() {
    if (!selectedAvatarFile) return;
    setUploadingAvatar(true);
    setImageError("");
    try {
      await uploadAvatar(selectedAvatarFile);
      clearAvatarSelection();
      setAvatarUrl("");
      setFailedAvatarUrl("");
      setProfileSuccess("Đã lưu avatar vào tài khoản.");
    } catch (err) {
      setImageError(err?.response?.data?.message || "Không thể tải avatar lên server.");
    } finally {
      setUploadingAvatar(false);
    }
  }

  async function handleVideoUpload() {
    if (!selectedVideoFile) return;
    setUploadingVideo(true);
    setVideoError("");
    try {
      await uploadProfileVideo(selectedVideoFile);
      clearVideoSelection();
      setProfileSuccess("Đã lưu video hồ sơ vào tài khoản.");
    } catch (err) {
      setVideoError(err?.response?.data?.message || "Không thể tải video lên server.");
    } finally {
      setUploadingVideo(false);
    }
  }

  const currentAvatarURL = resolveAssetUrl(user?.avatar_url);
  const currentProfileVideoURL = resolveAssetUrl(user?.profile_video_url);
  const avatarPreview = localImagePreview || (failedAvatarUrl === avatarUrl ? "" : avatarUrl.trim());
  const profileChanged =
    profileForm.name.trim() !== (user?.name || "").trim() || profileForm.phone.trim() !== (user?.phone || "").trim();

  return (
    <div className="grid profile-layout">
      <Card className="profile-card">
        <div className="page-header compact-header">
          <div>
            <span className="eyebrow">Tài khoản</span>
            <h1>Thông tin hồ sơ</h1>
            <p className="muted">Quản lý tên hiển thị, số điện thoại, avatar và video hồ sơ của tài khoản hiện tại.</p>
          </div>
          <Badge tone={user?.role === "admin" ? "primary" : "default"}>
            {user?.role === "admin" ? "Quản trị viên" : "Người dùng"}
          </Badge>
        </div>

        <div className="profile-readonly-grid">
          <div>
            <span>Họ tên</span>
            <strong>{user?.name || "—"}</strong>
          </div>
          <div>
            <span>Email</span>
            <strong>{user?.email || "—"}</strong>
          </div>
          <div>
            <span>Số điện thoại</span>
            <strong>{user?.phone || "Chưa cập nhật"}</strong>
          </div>
          <div>
            <span>Vai trò</span>
            <strong>{user?.role === "admin" ? "Quản trị viên" : "Người dùng"}</strong>
          </div>
          <div>
            <span>Avatar hiện tại</span>
            <strong>{user?.avatar_url ? "Đã có" : "Chưa có"}</strong>
          </div>
          <div>
            <span>Video hồ sơ hiện tại</span>
            <strong>{user?.profile_video_url ? "Đã có" : "Chưa có"}</strong>
          </div>
        </div>

        {currentAvatarURL ? (
          <div className="profile-media-preview profile-image-preview profile-avatar-current">
            <img src={currentAvatarURL} alt="Avatar tài khoản hiện tại" />
          </div>
        ) : null}

        {currentProfileVideoURL ? (
          <div className="profile-media-preview profile-video-preview profile-avatar-current">
            <video src={currentProfileVideoURL} controls playsInline aria-label="Video hồ sơ hiện tại" />
          </div>
        ) : null}

        <form className="profile-edit-form" onSubmit={handleProfileSubmit}>
          <div className="profile-edit-grid">
            <Input
              label="Tên hiển thị"
              value={profileForm.name}
              maxLength="100"
              onChange={(event) => {
                setProfileForm((current) => ({ ...current, name: event.target.value }));
                setProfileError("");
              }}
            />
            <Input
              label="Số điện thoại"
              type="tel"
              value={profileForm.phone}
              maxLength="20"
              placeholder="0xxxxxxxxx hoặc +84xxxxxxxxx"
              onChange={(event) => {
                setProfileForm((current) => ({ ...current, phone: event.target.value }));
                setProfileError("");
              }}
            />
          </div>
          <div className="actions">
            <Button type="submit" disabled={savingProfile || !profileChanged}>
              {savingProfile ? "Đang lưu..." : "Lưu thay đổi"}
            </Button>
          </div>
        </form>
        <ErrorMessage message={profileError} />
      </Card>

      <Card className="profile-card">
        <div className="page-header compact-header">
          <div>
            <span className="eyebrow">Bảo mật</span>
            <h2>Đổi mật khẩu</h2>
            <p className="muted">Xác nhận mật khẩu hiện tại trước khi đặt mật khẩu mới cho tài khoản.</p>
          </div>
          <Badge tone="warning">Bảo mật</Badge>
        </div>

        <form className="profile-password-form" onSubmit={handlePasswordSubmit}>
          <Input
            label="Mật khẩu hiện tại"
            type="password"
            value={passwordForm.currentPassword}
            autoComplete="current-password"
            onChange={(event) => {
              setPasswordForm((current) => ({ ...current, currentPassword: event.target.value }));
              setPasswordError("");
            }}
          />
          <Input
            label="Mật khẩu mới"
            type="password"
            value={passwordForm.newPassword}
            autoComplete="new-password"
            onChange={(event) => {
              setPasswordForm((current) => ({ ...current, newPassword: event.target.value }));
              setPasswordError("");
            }}
          />
          <Input
            label="Nhập lại mật khẩu mới"
            type="password"
            value={passwordForm.confirmPassword}
            autoComplete="new-password"
            onChange={(event) => {
              setPasswordForm((current) => ({ ...current, confirmPassword: event.target.value }));
              setPasswordError("");
            }}
          />
          <div className="actions">
            <Button type="submit" disabled={savingPassword}>
              {savingPassword ? "Đang đổi mật khẩu..." : "Đổi mật khẩu"}
            </Button>
          </div>
        </form>
        <ErrorMessage message={passwordError} />
      </Card>

      <Card className="profile-media-card">
        <div className="page-header compact-header">
          <div>
            <span className="eyebrow">Media hồ sơ</span>
            <h2>Tải ảnh và video lên server</h2>
            <p className="muted">File tải lên sẽ được lưu thật trên backend local và hiển thị lại sau khi đăng nhập lại.</p>
          </div>
          <Badge tone="success">Đã hỗ trợ lưu</Badge>
        </div>

        <div className="profile-media-grid">
          <section className="media-preview-section">
            <h3>Avatar</h3>
            <Input
              label="Ảnh từ URL để xem thử"
              type="url"
              placeholder="https://example.com/avatar.jpg"
              value={avatarUrl}
              onChange={(event) => {
                setAvatarUrl(event.target.value);
                setFailedAvatarUrl("");
              }}
            />
            <label className="field">
              <span>Hoặc chọn ảnh trên thiết bị — tối đa 5 MB</span>
              <input
                type="file"
                accept=".jpg,.jpeg,.png,.webp,.gif,.avif,image/jpeg,image/png,image/webp,image/gif,image/avif"
                onChange={handleImageFile}
              />
            </label>
            <ErrorMessage message={imageError} />

            <div className="profile-media-preview profile-image-preview">
              {avatarPreview ? (
                <img
                  src={avatarPreview}
                  alt="Ảnh đại diện xem thử"
                  onError={() => {
                    if (localImagePreview) {
                      setImageError("Không đọc được nội dung file ảnh đã chọn.");
                      clearAvatarSelection();
                    } else {
                      setFailedAvatarUrl(avatarUrl);
                    }
                  }}
                />
              ) : currentAvatarURL ? (
                <img src={currentAvatarURL} alt="Avatar hiện tại" />
              ) : (
                <span>Chưa chọn ảnh</span>
              )}
            </div>

            <div className="actions">
              <Button type="button" onClick={handleAvatarUpload} disabled={!selectedAvatarFile || uploadingAvatar}>
                {uploadingAvatar ? "Đang tải avatar..." : "Lưu avatar vào server"}
              </Button>
              {localImagePreview ? (
                <Button type="button" variant="secondary" onClick={clearAvatarSelection} disabled={uploadingAvatar}>
                  Xóa ảnh đã chọn
                </Button>
              ) : null}
            </div>
          </section>

          <section className="media-preview-section">
            <h3>Video hồ sơ</h3>
            <label className="field">
              <span>Video MP4/WebM — tối đa 10 giây, 20 MB</span>
              <input type="file" accept=".mp4,.webm,video/mp4,video/webm" onChange={handleVideoFile} />
            </label>
            <ErrorMessage message={videoError} />

            <div className="profile-media-preview profile-video-preview">
              {localVideoPreview ? (
                <video src={localVideoPreview} controls playsInline aria-label="Video hồ sơ xem thử" />
              ) : currentProfileVideoURL ? (
                <video src={currentProfileVideoURL} controls playsInline aria-label="Video hồ sơ hiện tại" />
              ) : (
                <span>Chưa chọn video</span>
              )}
            </div>

            <div className="actions">
              <Button type="button" onClick={handleVideoUpload} disabled={!selectedVideoFile || uploadingVideo}>
                {uploadingVideo ? "Đang tải video..." : "Lưu video vào server"}
              </Button>
              {localVideoPreview ? (
                <Button type="button" variant="secondary" onClick={clearVideoSelection} disabled={uploadingVideo}>
                  Xóa video đã chọn
                </Button>
              ) : null}
            </div>
          </section>
        </div>

        <div className="preview-only-banner" role="note">
          URL xem thử vẫn dùng được để kiểm tra nhanh, nhưng chỉ file đã upload mới được lưu thật vào tài khoản.
        </div>
      </Card>

      <Toast message={profileSuccess} tone="success" onDismiss={() => setProfileSuccess("")} />
    </div>
  );
}
