package model

import (
	"context"
	"fmt"
	"time"

	"github.com/QuantumNous/new-api/common"
	"gorm.io/gorm"
)

// ConversationArchive 归档表 - 用于存储旧对话
// 与 Conversation 表结构相同，但用于存储历史数据
type ConversationArchive struct {
	Id               int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId           int    `json:"user_id" gorm:"index;not null"`
	Username         string `json:"username" gorm:"index;not null;default:''"`
	ModelName        string `json:"model_name" gorm:"index;not null;default:''"`
	TokenId          int    `json:"token_id" gorm:"index;default:0"`
	TokenName        string `json:"token_name" gorm:"default:''"`
	ChannelId        int    `json:"channel_id" gorm:"index;default:0"`
	RequestMessages  string `json:"request_messages" gorm:"type:text"`
	ResponseContent  string `json:"response_content" gorm:"type:text"`
	PromptTokens     int    `json:"prompt_tokens" gorm:"default:0"`
	CompletionTokens int    `json:"completion_tokens" gorm:"default:0"`
	TotalTokens      int    `json:"total_tokens" gorm:"default:0"`
	IsStream         bool   `json:"is_stream" gorm:"default:false"`
	CreatedAt        int64  `json:"created_at" gorm:"bigint;index;not null"`
	UseTime          int    `json:"use_time" gorm:"default:0"`
	Ip               string `json:"ip" gorm:"index;default:''"`
	Group            string `json:"group" gorm:"index;default:''"`
	ArchivedAt       int64  `json:"archived_at" gorm:"bigint;index"` // 归档时间
}

func (ConversationArchive) TableName() string {
	return "conversations_archive"
}

// ArchiveOldConversations 归档旧对话到归档表
// targetTimestamp: 归档此时间点之前的数据
// batchSize: 每批处理的记录数
func ArchiveOldConversations(ctx context.Context, targetTimestamp int64, batchSize int) (int64, error) {
	var totalArchived int64 = 0
	archivedAt := time.Now().Unix()

	for {
		if ctx.Err() != nil {
			return totalArchived, ctx.Err()
		}

		// 开启事务
		err := LOG_DB.Transaction(func(tx *gorm.DB) error {
			// 1. 查询需要归档的数据
			var conversations []Conversation
			if err := tx.Where("created_at < ?", targetTimestamp).
				Limit(batchSize).
				Find(&conversations).Error; err != nil {
				return err
			}

			if len(conversations) == 0 {
				return nil // 没有数据需要归档
			}

			// 2. 转换为归档记录
			archives := make([]ConversationArchive, len(conversations))
			ids := make([]int, len(conversations))
			for i, conv := range conversations {
				archives[i] = ConversationArchive{
					UserId:           conv.UserId,
					Username:         conv.Username,
					ModelName:        conv.ModelName,
					TokenId:          conv.TokenId,
					TokenName:        conv.TokenName,
					ChannelId:        conv.ChannelId,
					RequestMessages:  conv.RequestMessages,
					ResponseContent:  conv.ResponseContent,
					PromptTokens:     conv.PromptTokens,
					CompletionTokens: conv.CompletionTokens,
					TotalTokens:      conv.TotalTokens,
					IsStream:         conv.IsStream,
					CreatedAt:        conv.CreatedAt,
					UseTime:          conv.UseTime,
					Ip:               conv.Ip,
					Group:            conv.Group,
					ArchivedAt:       archivedAt,
				}
				ids[i] = conv.Id
			}

			// 3. 插入到归档表
			if err := tx.Create(&archives).Error; err != nil {
				return err
			}

			// 4. 从主表删除
			if err := tx.Where("id IN ?", ids).Delete(&Conversation{}).Error; err != nil {
				return err
			}

			totalArchived += int64(len(conversations))
			return nil
		})

		if err != nil {
			return totalArchived, err
		}

		// 如果处理的数据少于 batchSize，说明已经处理完了
		if totalArchived%int64(batchSize) != 0 {
			break
		}

		// 休眠一下，避免持续占用数据库资源
		time.Sleep(100 * time.Millisecond)
	}

	return totalArchived, nil
}

// GetConversationsWithArchive 从主表和归档表查询对话（统一查询接口）
func GetConversationsWithArchive(userId int, modelName string, username string, startTime int64, endTime int64, startIdx int, num int, searchArchive bool) ([]*Conversation, int64, error) {
	// 默认只查询主表
	conversations, total, err := GetConversations(userId, modelName, username, startTime, endTime, startIdx, num)
	if err != nil {
		return nil, 0, err
	}

	// 如果需要查询归档表
	if searchArchive {
		archives, archiveTotal, err := GetArchivedConversations(userId, modelName, username, startTime, endTime, startIdx, num)
		if err != nil {
			return conversations, total, err // 归档查询失败不影响主表结果
		}

		// 合并结果
		for _, archive := range archives {
			conv := &Conversation{
				Id:               archive.Id,
				UserId:           archive.UserId,
				Username:         archive.Username,
				ModelName:        archive.ModelName,
				TokenId:          archive.TokenId,
				TokenName:        archive.TokenName,
				ChannelId:        archive.ChannelId,
				RequestMessages:  archive.RequestMessages,
				ResponseContent:  archive.ResponseContent,
				PromptTokens:     archive.PromptTokens,
				CompletionTokens: archive.CompletionTokens,
				TotalTokens:      archive.TotalTokens,
				IsStream:         archive.IsStream,
				CreatedAt:        archive.CreatedAt,
				UseTime:          archive.UseTime,
				Ip:               archive.Ip,
				Group:            archive.Group,
			}
			conversations = append(conversations, conv)
		}
		total += archiveTotal
	}

	return conversations, total, nil
}

// GetArchivedConversations 查询归档表
func GetArchivedConversations(userId int, modelName string, username string, startTime int64, endTime int64, startIdx int, num int) ([]*ConversationArchive, int64, error) {
	var archives []*ConversationArchive
	var total int64

	tx := LOG_DB.Model(&ConversationArchive{})

	// 筛选条件
	if userId > 0 {
		tx = tx.Where("user_id = ?", userId)
	}
	if modelName != "" {
		tx = tx.Where("model_name LIKE ?", "%"+modelName+"%")
	}
	if username != "" {
		tx = tx.Where("username = ?", username)
	}
	if startTime > 0 {
		tx = tx.Where("created_at >= ?", startTime)
	}
	if endTime > 0 {
		tx = tx.Where("created_at <= ?", endTime)
	}

	// 获取总数
	err := tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	err = tx.Order("created_at DESC").Limit(num).Offset(startIdx).Find(&archives).Error
	return archives, total, err
}

// CleanupOldArchives 清理归档表中的超旧数据（比如1年以前的）
func CleanupOldArchives(ctx context.Context, targetTimestamp int64, batchSize int) (int64, error) {
	var totalDeleted int64 = 0

	for {
		if ctx.Err() != nil {
			return totalDeleted, ctx.Err()
		}

		result := LOG_DB.Where("created_at < ?", targetTimestamp).
			Limit(batchSize).
			Delete(&ConversationArchive{})

		if result.Error != nil {
			return totalDeleted, result.Error
		}

		totalDeleted += result.RowsAffected

		if result.RowsAffected < int64(batchSize) {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	return totalDeleted, nil
}

// GetTableSize 获取表的大小信息（仅支持 MySQL 和 PostgreSQL）
func GetTableSize(tableName string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if common.UsingPostgreSQL {
		// PostgreSQL
		var tableSize int64
		var indexSize int64
		err := LOG_DB.Raw(`
			SELECT
				pg_total_relation_size(?) as total_size,
				pg_indexes_size(?) as index_size
		`, tableName, tableName).Row().Scan(&tableSize, &indexSize)

		if err != nil {
			return nil, err
		}

		result["table_size"] = tableSize
		result["index_size"] = indexSize
		result["total_size"] = tableSize + indexSize
	} else if !common.UsingSQLite {
		// MySQL
		var dataLength int64
		var indexLength int64
		err := LOG_DB.Raw(`
			SELECT
				data_length as data_length,
				index_length as index_length
			FROM information_schema.TABLES
			WHERE table_schema = DATABASE()
			AND table_name = ?
		`, tableName).Row().Scan(&dataLength, &indexLength)

		if err != nil {
			return nil, err
		}

		result["table_size"] = dataLength
		result["index_size"] = indexLength
		result["total_size"] = dataLength + indexLength
	} else {
		// SQLite 不支持
		return nil, fmt.Errorf("SQLite does not support table size query")
	}

	return result, nil
}

// OptimizeConversationTable 优化对话表（重建索引、回收空间）
func OptimizeConversationTable() error {
	if common.UsingPostgreSQL {
		// PostgreSQL: VACUUM ANALYZE
		if err := LOG_DB.Exec("VACUUM ANALYZE conversations").Error; err != nil {
			return err
		}
		if err := LOG_DB.Exec("VACUUM ANALYZE conversations_archive").Error; err != nil {
			return err
		}
	} else if !common.UsingSQLite {
		// MySQL: OPTIMIZE TABLE
		if err := LOG_DB.Exec("OPTIMIZE TABLE conversations").Error; err != nil {
			return err
		}
		if err := LOG_DB.Exec("OPTIMIZE TABLE conversations_archive").Error; err != nil {
			return err
		}
	}

	return nil
}

// GetConversationTableStats 获取对话表统计信息
func GetConversationTableStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 主表统计
	var mainCount int64
	if err := LOG_DB.Model(&Conversation{}).Count(&mainCount).Error; err != nil {
		return nil, err
	}
	stats["main_table_count"] = mainCount

	// 归档表统计
	var archiveCount int64
	if err := LOG_DB.Model(&ConversationArchive{}).Count(&archiveCount).Error; err != nil {
		stats["archive_table_count"] = 0 // 归档表可能不存在
	} else {
		stats["archive_table_count"] = archiveCount
	}

	// 总计
	stats["total_count"] = mainCount + archiveCount

	// 表大小（仅 MySQL/PostgreSQL）
	if !common.UsingSQLite {
		mainSize, err := GetTableSize("conversations")
		if err == nil {
			stats["main_table_size"] = mainSize
		}

		archiveSize, err := GetTableSize("conversations_archive")
		if err == nil {
			stats["archive_table_size"] = archiveSize
		}
	}

	// 最老和最新的记录
	var oldestTime, newestTime int64
	LOG_DB.Model(&Conversation{}).Select("MIN(created_at)").Scan(&oldestTime)
	LOG_DB.Model(&Conversation{}).Select("MAX(created_at)").Scan(&newestTime)
	stats["oldest_conversation"] = oldestTime
	stats["newest_conversation"] = newestTime

	return stats, nil
}
